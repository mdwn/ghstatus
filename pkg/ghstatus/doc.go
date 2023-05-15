// Package ghstatus contains the Github status client.
//
// This package contains the Github Status client and a few convience
// methods for rendering them. The client parses the raw JSON from the
// Github Status API, which is documented (incompletely) here:
// https://www.githubstatus.com/api
//
// This is a relatively simple API. Retries were added by using a
// retryable HTTP client, but no other features pertaining to API rate limiting
// or caching have been added. As of this writing, no API rate limits are documented
// in the API, so it is currently assumed there either is no rate limit or it is very high.
package ghstatus
