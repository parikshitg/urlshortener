# CORS Configuration

This URL shortener service includes CORS (Cross-Origin Resource Sharing) support to allow web applications from different domains to access the API.

## Default Configuration

By default, the service uses a permissive CORS configuration suitable for development:

- **Allowed Origins**: `*` (all origins)
- **Allowed Methods**: `GET`, `POST`, `PUT`, `DELETE`, `OPTIONS`
- **Allowed Headers**: `*` (all headers)
- **Exposed Headers**: `Content-Length`
- **Allow Credentials**: `false`
- **Max Age**: `43200` seconds (12 hours)

## Environment Variables

You can configure CORS behavior using the following environment variables:

### CORS_ALLOWED_ORIGINS
Comma-separated list of allowed origins. Default: `*`

```bash
# Allow all origins (default)
CORS_ALLOWED_ORIGINS="*"

# Allow specific origins
CORS_ALLOWED_ORIGINS="https://example.com,https://app.example.com"
```

### CORS_ALLOWED_METHODS
Comma-separated list of allowed HTTP methods. Default: `GET,POST,PUT,DELETE,OPTIONS`

```bash
# Allow specific methods
CORS_ALLOWED_METHODS="GET,POST,OPTIONS"
```

### CORS_ALLOWED_HEADERS
Comma-separated list of allowed headers. Default: `*`

```bash
# Allow all headers (default)
CORS_ALLOWED_HEADERS="*"

# Allow specific headers
CORS_ALLOWED_HEADERS="Content-Type,Authorization"
```

### CORS_MAX_AGE
Maximum age for preflight requests in seconds. Default: `43200` (12 hours)

```bash
# Set max age to 1 hour
CORS_MAX_AGE="3600"
```

### CORS_ALLOW_CREDENTIALS
Whether to allow credentials in CORS requests. Default: `false`

```bash
# Allow credentials
CORS_ALLOW_CREDENTIALS="true"
```

## Production Configuration Example

For production, you should use a more restrictive configuration:

```bash
# Only allow specific origins
CORS_ALLOWED_ORIGINS="https://yourdomain.com,https://app.yourdomain.com"

# Only allow necessary methods
CORS_ALLOWED_METHODS="GET,POST,OPTIONS"

# Only allow necessary headers
CORS_ALLOWED_HEADERS="Content-Type,Authorization"

# Shorter cache time for preflight requests
CORS_MAX_AGE="3600"

# Allow credentials if needed
CORS_ALLOW_CREDENTIALS="true"
```

## Testing CORS

You can test CORS functionality using curl:

```bash
# Test preflight request
curl -X OPTIONS \
  -H "Origin: https://example.com" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type" \
  http://localhost:8080/v1/shorten

# Test actual request
curl -X POST \
  -H "Origin: https://example.com" \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com"}' \
  http://localhost:8080/v1/shorten
```

## Browser Testing

You can test CORS from a browser console:

```javascript
// Test from a different origin
fetch('http://localhost:8080/health/', {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json',
  },
})
.then(response => response.json())
.then(data => console.log(data))
.catch(error => console.error('CORS Error:', error));
```

## Security Considerations

1. **Never use `*` for origins in production** - Always specify exact domains
2. **Be restrictive with methods** - Only allow the HTTP methods you actually use
3. **Limit headers** - Only allow necessary headers
4. **Use HTTPS** - Always use HTTPS in production
5. **Monitor CORS errors** - Log and monitor CORS-related errors

## Troubleshooting

### Common CORS Issues

1. **"Access to fetch at '...' from origin '...' has been blocked by CORS policy"**
   - Check that your origin is in the `CORS_ALLOWED_ORIGINS` list
   - Verify the request method is allowed in `CORS_ALLOWED_METHODS`

2. **Preflight requests failing**
   - Ensure `OPTIONS` method is in `CORS_ALLOWED_METHODS`
   - Check that required headers are in `CORS_ALLOWED_HEADERS`

3. **Credentials not working**
   - Set `CORS_ALLOW_CREDENTIALS="true"`
   - Cannot use `*` for origins when credentials are enabled

### Debugging

Enable debug logging to see CORS-related information:

```bash
LOG_LEVEL="debug" ./urlshortener
```

This will show detailed information about CORS headers being set and requests being processed.
