// Package signv4 provides an high level implementation of AWS Signature V4 protocol
// built around AWS SDK aws/signer/v4 package.
//
// Features provided by this package are:
// - Computation of the body hash (SHA-256);
// - Automatic credentials renewal when necessary;
package signv4

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	signerv4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/morelj/httptools/request"
)

const (
	awsService  = "execute-api"
	minValidity = 15 * time.Second
)

// emptySha256 is the SHA-256 of an empty value
const emptySha256 = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

// A TimeProvider is a function which returns the time to use in the signature.
// The most useful TimeProvider is time.Now which returns the current time.
type TimeProvider func() time.Time

// FixedTime returns a TimeProvider which always return the time t.
// It is mainly used to test signatures.
func FixedTime(t time.Time) TimeProvider {
	return func() time.Time {
		return t
	}
}

type Signer struct {
	signer       *signerv4.Signer
	provider     aws.CredentialsProvider
	region       string
	service      string
	timeProvider TimeProvider
	minValidity  time.Duration
	credentials  chan aws.Credentials
}

// New returns a new Signer which will use the given provider to provide credentials.
// If provider returns credentials which can expire, new credentials are obtained each time the previous
// credentials are about to expire.
//
// By default, the time used to sign is the current time (at the moment of the signature) and the service name
// is "execute-api".
//
// By default, new credentials are obtained when previous ones expire in less than 15 seconds.
//
// These parameters can be modified by passing Options.
//
// A single Signer can be used concurrently by multiple goroutines.
func New(provider aws.CredentialsProvider, region string, options ...Option) *Signer {
	s := &Signer{
		signer:       signerv4.NewSigner(),
		provider:     provider,
		region:       region,
		service:      awsService,
		timeProvider: time.Now,
		minValidity:  minValidity,
		credentials:  make(chan aws.Credentials, 1), // The buffer of 1 is required
	}

	for _, option := range options {
		option(s)
	}

	// Put initial invalid credentials so getCredentials won't wait forever on the first call
	s.credentials <- aws.Credentials{}

	return s
}

// SignRequest signs an [*request.Request].
// If req.Err is not nil, the request is not signed.
// req.Err is set to any error which may occur during the signing process.
func (s *Signer) SignRequest(ctx context.Context, req *request.Request) {
	if req.Err == nil {
		req.Err = s.SignHTTPRequest(ctx, req.Request)
	}
}

// SignHTTPRequest signs an [*http.Request]
// A non-nil error is returned if the signature fails.
func (s *Signer) SignHTTPRequest(ctx context.Context, req *http.Request) error {
	credentials, err := s.getCredentials(ctx)
	if err != nil {
		return err
	}

	payloadHash, err := getBodyHash(req)
	if err != nil {
		return fmt.Errorf("failed to compute payload hash: %w", err)
	}

	return s.signer.SignHTTP(ctx, credentials, req, payloadHash, s.service, s.region, s.timeProvider())
}

// getCredentials return valid credentials.
// Credentials are first obtained from the s.credentials channel. If these credentials are invalid, new one are
// retrieved using the provider.
// Credentials are put back in s.credentials before returning from this function, to make them available to other
// goroutines.
func (s *Signer) getCredentials(ctx context.Context) (aws.Credentials, error) {
	// Obtain current credentials
	select {
	case credentials := <-s.credentials:
		var err error
		if !s.isValid(&credentials) {
			// Retrieve new credentials
			credentials, err = s.provider.Retrieve(ctx)
			if err != nil {
				credentials = aws.Credentials{}
				err = fmt.Errorf("failed to retrieve AWS credentials: %w", err)
			}
		}

		// Put the credentials to the channel so the next call can use them immediately
		s.credentials <- credentials
		return credentials, err

	case <-ctx.Done():
		return aws.Credentials{}, fmt.Errorf("context timeout while waiting for credentials")
	}
}

func (s *Signer) isValid(credentials *aws.Credentials) bool {
	if credentials == nil || credentials.AccessKeyID == "" {
		// No credentials defined
		return false
	}

	if credentials.CanExpire && credentials.Expires.Before(time.Now().Add(s.minValidity)) {
		// Credentials are expired
		return false
	}

	return true
}

func getBodyHash(req *http.Request) (string, error) {
	if req.Body == nil || req.Body == http.NoBody {
		return emptySha256, nil
	}

	if req.GetBody == nil {
		return "", errors.New("cannot get body of request")
	}
	body, err := req.GetBody()
	if err != nil {
		return "", err
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, body); err != nil {
		return "", fmt.Errorf("failed to compute body SHA-256: %w", err)
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// simpleCredentialsProvider provides a simple aws.CredentialsProvider which returns itself each time the Retrieve()
// method is called.
type simpleCredentialsProvider aws.Credentials

func (scp simpleCredentialsProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials(scp), nil
}

// An Option applies some customization to a Signer.
type Option func(s *Signer)

// WithService returns an Option which will override the service used in signatures with the provided value.
func WithService(service string) Option {
	return func(s *Signer) {
		s.service = service
	}
}

// WithTimeProvider returns an Option which will override the TimeProvider used in signatures with the provided one.
func WithTimeProvider(timeProvider TimeProvider) Option {
	return func(s *Signer) {
		s.timeProvider = timeProvider
	}
}

// WithMinValidity returns an Option which will override the time.Duration
func WithMinValidity(minValidity time.Duration) Option {
	return func(s *Signer) {
		s.minValidity = minValidity
	}
}
