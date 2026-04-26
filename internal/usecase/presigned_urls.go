package usecase

func (s *MapService) rewritePresignedUploadURL(rawURL string) (string, error) {
	return s.rewritePresignedURL(rawURL, s.uploadBaseProxyURL)
}

func (s *MapService) rewritePresignedDownloadURL(rawURL string) (string, error) {
	return s.rewritePresignedURL(rawURL, s.downloadBaseProxyURL)
}

func (s *MapService) rewritePresignedURL(rawURL, relayBaseURL string) (string, error) {
	if !s.proxyEnabled {
		return rawURL, nil
	}

	return RewriteToRelay(rawURL, relayBaseURL)
}
