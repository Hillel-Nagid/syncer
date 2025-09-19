# Backend Testing Plan

This document outlines a comprehensive testing strategy for the Syncer backend application. The testing plan covers unit tests, integration tests, API tests, and security tests for both implemented and planned features.

## Table of Contents

1. [Testing Strategy Overview](#testing-strategy-overview)
2. [Core Components Testing](#core-components-testing)
3. [API Layer Testing](#api-layer-testing)
4. [Services Layer Testing](#services-layer-testing)
5. [Database Layer Testing](#database-layer-testing)
6. [Security Testing](#security-testing)
7. [Unimplemented Features Testing](#unimplemented-features-testing)
8. [Testing Infrastructure](#testing-infrastructure)
9. [CI/CD Integration](#cicd-integration)
10. [Test Data Management](#test-data-management)

## Testing Strategy Overview

### Testing Principles

1. **Test-Driven Development (TDD)**: Write tests before implementing new features
2. **Privacy-First Testing**: Ensure no user data leaks in tests
3. **Comprehensive Coverage**: Unit, integration, and end-to-end tests
4. **Mock External Dependencies**: Use mocks for external APIs
5. **Security Testing**: Include security and vulnerability tests
6. **Performance Testing**: Include load and performance tests

### Testing Pyramid

```
End-to-End Tests (5%)
Integration Tests (20%)
Unit Tests (75%)
```

### Test Categories

- **Unit Tests**: Individual functions, methods, and components
- **Integration Tests**: Component interactions and external service integrations
- **API Tests**: REST endpoint testing with authentication
- **Security Tests**: Authentication, authorization, encryption
- **Database Tests**: Migrations, queries, transactions
- **Performance Tests**: Load testing and benchmarking

## Core Components Testing

### 1. Authentication System Testing

#### JWT Service (`backend/core/auth/jwt.go`)

```go
// backend/core/auth/jwt_test.go
func TestJWTService_GenerateAccessToken(t *testing.T) {
    service := NewJWTService("test-secret")
    userID := "user123"

    token, err := service.GenerateAccessToken(userID)
    assert.NoError(t, err)
    assert.NotEmpty(t, token)

    // Verify token structure
    claims, err := service.ValidateToken(token)
    assert.NoError(t, err)
    assert.Equal(t, userID, claims.UserID)
}

func TestJWTService_TokenExpiration(t *testing.T) {
    service := NewJWTService("test-secret")

    // Test expired token
    expiredToken := generateExpiredToken(t)
    _, err := service.ValidateToken(expiredToken)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "expired")
}

func TestJWTService_InvalidToken(t *testing.T) {
    service := NewJWTService("test-secret")

    invalidTokens := []string{
        "invalid-token",
        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", // Incomplete
        "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoiMTIzIn0", // No signature
    }

    for _, token := range invalidTokens {
        _, err := service.ValidateToken(token)
        assert.Error(t, err, "Token should be invalid: %s", token)
    }
}
```

#### Session Service (`backend/core/auth/session.go`)

```go
// backend/core/auth/session_test.go
func TestSessionService_CreateSession(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    service := NewSessionService(db)
    userID := "user123"
    ip := "127.0.0.1"
    userAgent := "test-agent"

    session, err := service.CreateSession(userID, ip, userAgent)
    assert.NoError(t, err)
    assert.NotEmpty(t, session.ID)
    assert.Equal(t, userID, session.UserID)
    assert.Equal(t, ip, session.IPAddress)
    assert.Equal(t, userAgent, session.UserAgent)
}

func TestSessionService_ValidateSession(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    service := NewSessionService(db)

    // Test valid session
    session := createTestSession(t, db)
    valid, err := service.ValidateSession(session.SessionToken)
    assert.NoError(t, err)
    assert.True(t, valid)

    // Test invalid session
    valid, err = service.ValidateSession("invalid-token")
    assert.NoError(t, err)
    assert.False(t, valid)

    // Test expired session
    expiredSession := createExpiredSession(t, db)
    valid, err = service.ValidateSession(expiredSession.SessionToken)
    assert.NoError(t, err)
    assert.False(t, valid)
}

func TestSessionService_RefreshSession(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    service := NewSessionService(db)
    session := createTestSession(t, db)

    newSession, err := service.RefreshSession(session.RefreshToken)
    assert.NoError(t, err)
    assert.NotEmpty(t, newSession.AccessToken)
    assert.NotEmpty(t, newSession.RefreshToken)
    assert.NotEqual(t, session.AccessToken, newSession.AccessToken)
    assert.NotEqual(t, session.RefreshToken, newSession.RefreshToken)
}

func TestSessionService_RevokeSession(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    service := NewSessionService(db)
    session := createTestSession(t, db)

    err := service.RevokeSession(session.ID)
    assert.NoError(t, err)

    // Verify session is revoked
    valid, err := service.ValidateSession(session.SessionToken)
    assert.NoError(t, err)
    assert.False(t, valid)
}
```

### 2. Services Layer Testing

#### Service Registry (`backend/core/services/registry.go`)

```go
// backend/core/services/registry_test.go
func TestServiceRegistry_Register(t *testing.T) {
    registry := NewServiceRegistry(nil, slog.Default())
    mockService := &MockService{name: "test-service"}

    err := registry.Register(mockService)
    assert.NoError(t, err)

    service, err := registry.GetService("test-service")
    assert.NoError(t, err)
    assert.Equal(t, mockService, service)
}

func TestServiceRegistry_DuplicateRegistration(t *testing.T) {
    registry := NewServiceRegistry(nil, slog.Default())
    mockService := &MockService{name: "test-service"}

    err := registry.Register(mockService)
    assert.NoError(t, err)

    err = registry.Register(mockService)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "already registered")
}

func TestServiceRegistry_GetNonExistentService(t *testing.T) {
    registry := NewServiceRegistry(nil, slog.Default())

    service, err := registry.GetService("non-existent")
    assert.Error(t, err)
    assert.Nil(t, service)
    assert.Contains(t, err.Error(), "not found")
}

func TestServiceRegistry_ThreadSafety(t *testing.T) {
    registry := NewServiceRegistry(nil, slog.Default())

    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            service := &MockService{name: fmt.Sprintf("service-%d", id)}
            err := registry.Register(service)
            assert.NoError(t, err)
        }(i)
    }

    wg.Wait()

    services := registry.ListServices()
    assert.Len(t, services, 10)
}
```

#### OAuth Manager (`backend/core/services/oauth.go`)

```go
// backend/core/services/oauth_test.go
func TestOAuthManager_InitiateAuth(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    registry := setupTestRegistry(t)
    manager := NewOAuthManager(registry, db, []byte("test-encryption-key"))

    auth, err := manager.InitiateAuth("spotify", "user123", "http://localhost:3000/callback")
    assert.NoError(t, err)
    assert.NotEmpty(t, auth.AuthURL)
    assert.NotEmpty(t, auth.State)

    // Verify pending auth was stored
    storedAuth := getPendingAuth(t, db, auth.State)
    assert.Equal(t, "user123", storedAuth.UserID)
    assert.Equal(t, "spotify", storedAuth.ServiceName)
}

func TestOAuthManager_TokenEncryption(t *testing.T) {
    encryption := security.NewTokenEncryption([]byte("test-encryption-key"))

    accessToken := "test-access-token-12345"
    refreshToken := "test-refresh-token-67890"

    // Encrypt tokens
    encryptedAccess, encryptedRefresh, err := encryption.EncryptTokens(accessToken, refreshToken)
    assert.NoError(t, err)

    // Verify encrypted tokens are different from originals
    assert.NotEqual(t, accessToken, encryptedAccess)
    assert.NotEqual(t, refreshToken, encryptedRefresh)

    // Decrypt tokens
    decryptedAccess, decryptedRefresh, err := encryption.DecryptTokens(encryptedAccess, encryptedRefresh)
    assert.NoError(t, err)

    // Verify decryption works
    assert.Equal(t, accessToken, decryptedAccess)
    assert.Equal(t, refreshToken, decryptedRefresh)

    // Test with wrong key
    wrongEncryption := security.NewTokenEncryption([]byte("wrong-key"))
    _, _, err = wrongEncryption.DecryptTokens(encryptedAccess, encryptedRefresh)
    assert.Error(t, err)
}
```

### 3. Sync Engine Testing

#### Sync Engine (`backend/core/sync/engine.go`)

```go
// backend/core/sync/engine_test.go
func TestSyncEngine_QueueManualSync(t *testing.T) {
    registry := setupTestRegistry(t)
    oauth := setupTestOAuth(t)
    db := setupTestDB(t)
    defer db.Close()

    engine := NewSyncEngine(registry, oauth, db, 2)

    req := &SyncJobRequest{
        UserID: "user123",
        ServicePairs: []ServicePair{{
            SourceService: "spotify",
            TargetService: "deezer",
            SyncMode:      SyncModeFrom,
        }},
        SyncType:    SyncTypeFavorites,
        RequestedAt: time.Now(),
    }

    err := engine.QueueManualSync(req)
    assert.NoError(t, err)
}

func TestSyncEngine_ValidateSyncRequest(t *testing.T) {
    engine := NewSyncEngine(nil, nil, nil, 2)

    tests := []struct {
        name        string
        req         *SyncJobRequest
        shouldError bool
        errorMsg    string
    }{
        {
            name: "valid request",
            req: &SyncJobRequest{
                ServicePairs: []ServicePair{{
                    SourceService: "spotify",
                    TargetService: "deezer",
                    SyncMode:      SyncModeFrom,
                }},
            },
            shouldError: false,
        },
        {
            name:        "no service pairs",
            req:         &SyncJobRequest{ServicePairs: []ServicePair{}},
            shouldError: true,
            errorMsg:    "at least one service pair is required",
        },
        {
            name: "same source and target",
            req: &SyncJobRequest{
                ServicePairs: []ServicePair{{
                    SourceService: "spotify",
                    TargetService: "spotify",
                    SyncMode:      SyncModeFrom,
                }},
            },
            shouldError: true,
            errorMsg:    "source and target services must be different",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := engine.validateSyncRequest(tt.req)
            if tt.shouldError {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errorMsg)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

#### Track Transformer (`backend/core/sync/transformer.go`)

```go
// backend/core/sync/transformer_test.go
func TestTrackTransformer_SpotifyToUniversal(t *testing.T) {
    transformer := NewTrackTransformer()

    spotifyTrack := SpotifyTrack{
        ID:        "spotify123",
        Name:      "Test Track",
        Artists:   []SpotifyArtist{{Name: "Test Artist"}},
        Album:     SpotifyAlbum{Name: "Test Album"},
        Duration:  180000, // 3 minutes in ms
        ExternalIDs: map[string]string{"isrc": "US1234567890"},
    }

    universal := transformer.SpotifyToUniversal(spotifyTrack)

    assert.Equal(t, "Test Track", universal.Title)
    assert.Equal(t, "Test Artist", universal.Artist)
    assert.Equal(t, "Test Album", universal.Album)
    assert.Equal(t, 180000, universal.Duration)
    assert.Equal(t, "US1234567890", universal.ISRC)
    assert.Equal(t, "spotify123", universal.ExternalIDs["spotify"])
}

func TestTrackTransformer_FindBestMatch(t *testing.T) {
    transformer := NewTrackTransformer()

    sourceTrack := UniversalTrack{
        Title:    "Test Track",
        Artist:   "Test Artist",
        Album:    "Test Album",
        Duration: 180000,
        ISRC:     "US1234567890",
    }

    candidates := []UniversalTrack{
        {
            Title:    "Test Track",
            Artist:   "Test Artist",
            Album:    "Test Album",
            Duration: 180000,
            ISRC:     "US1234567890",
        },
        {
            Title:    "Different Track",
            Artist:   "Different Artist",
            Album:    "Different Album",
            Duration: 120000,
        },
    }

    bestMatch, score := transformer.FindBestMatch(sourceTrack, candidates, 0.5)

    assert.NotNil(t, bestMatch)
    assert.Equal(t, 1.0, score) // Perfect match due to ISRC
    assert.Equal(t, "Test Track", bestMatch.Title)
}
```

## API Layer Testing

### 1. Authentication Endpoints

```go
// backend/api/auth_test.go
func TestAuthRegister(t *testing.T) {
    router := setupTestRouter(t)
    db := getTestDB(t)

    payload := map[string]interface{}{
        "email":      "test@example.com",
        "password":   "testpassword123",
        "full_name":  "Test User",
    }

    req, _ := http.NewRequest("POST", "/auth/register", toJSONReader(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)

    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)

    // Verify user was created
    var user User
    err = db.Get(&user, "SELECT * FROM users WHERE email = $1", payload["email"])
    assert.NoError(t, err)
    assert.Equal(t, payload["email"], user.Email)
    assert.Equal(t, payload["full_name"], user.FullName)
}

func TestAuthProtectedEndpoint(t *testing.T) {
    router := setupTestRouter(t)
    db := getTestDB(t)

    user := createTestUser(t, db, "test@example.com", "hashedpassword", "Test User")
    session := createTestSession(t, db, user.ID)

    // Test without authentication
    req, _ := http.NewRequest("GET", "/auth/profile", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    assert.Equal(t, 401, w.Code)

    // Test with valid session cookie
    req, _ = http.NewRequest("GET", "/auth/profile", nil)
    req.AddCookie(&http.Cookie{
        Name:  "access_token",
        Value: session.AccessToken,
    })

    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)
    assert.Equal(t, 200, w.Code)
}
```

### 2. Sync Endpoints

```go
// backend/api/sync_test.go
func TestInitiateManualSync(t *testing.T) {
    router := setupTestRouter(t)
    db := getTestDB(t)

    user := createTestUser(t, db, "test@example.com", "hashedpassword", "Test User")
    session := createTestSession(t, db, user.ID)

    // Create connected services
    createConnectedService(t, db, user.ID, "spotify")
    createConnectedService(t, db, user.ID, "deezer")

    payload := map[string]interface{}{
        "service_pairs": []map[string]interface{}{
            {
                "source_service": "spotify",
                "target_service": "deezer",
                "sync_mode":      "sync-from",
            },
        },
        "sync_type": "favorites",
        "sync_options": map[string]interface{}{
            "dry_run": false,
        },
    }

    req, _ := http.NewRequest("POST", "/api/sync/manual", toJSONReader(payload))
    req.Header.Set("Content-Type", "application/json")
    req.AddCookie(&http.Cookie{
        Name:  "access_token",
        Value: session.AccessToken,
    })

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
}
```

## Services Layer Testing

### 1. Spotify Service Testing

```go
// backend/services/music/spotify/service_test.go
func TestSpotifyService_GetAuthURL(t *testing.T) {
    service := NewSpotifyService()

    state := "test-state"
    redirectURL := "http://localhost:3000/callback"

    authURL, err := service.GetAuthURL(state, redirectURL)
    assert.NoError(t, err)
    assert.Contains(t, authURL, "https://accounts.spotify.com/authorize")
    assert.Contains(t, authURL, "client_id="+service.clientID)
    assert.Contains(t, authURL, "state="+state)
}

func TestSpotifyService_SyncUserData(t *testing.T) {
    service := NewSpotifyService()

    // Mock HTTP server
    server := setupMockSpotifyServer(t)
    defer server.Close()

    tokens := &OAuthTokens{
        AccessToken: "test-access-token",
        TokenType:   "Bearer",
        ExpiresAt:   time.Now().Add(time.Hour),
    }

    result, err := service.SyncUserData(context.Background(), tokens, time.Now().Add(-24*time.Hour))
    assert.NoError(t, err)
    assert.True(t, result.Success)
    assert.Greater(t, len(result.Items), 0)
}
```

## Security Testing

### 1. Authentication Security

```go
// backend/core/auth/security_test.go
func TestJWTService_TokenTampering(t *testing.T) {
    service := NewJWTService("test-secret")

    // Generate valid token
    originalToken, err := service.GenerateAccessToken("user123")
    assert.NoError(t, err)

    // Tamper with token (modify payload)
    parts := strings.Split(originalToken, ".")
    payload, err := base64.RawURLEncoding.DecodeString(parts[1])
    assert.NoError(t, err)

    var claims map[string]interface{}
    err = json.Unmarshal(payload, &claims)
    assert.NoError(t, err)

    // Modify user ID
    claims["user_id"] = "hacked-user"
    tamperedPayload, err := json.Marshal(claims)
    assert.NoError(t, err)

    tamperedToken := parts[0] + "." + base64.RawURLEncoding.EncodeToString(tamperedPayload) + "." + parts[2]

    // Verify tampered token is rejected
    _, err = service.ValidateToken(tamperedToken)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "invalid signature")
}
```

## Unimplemented Features Testing

### 1. Apple Music Service (Planned)

```go
// backend/services/music/apple/service_test.go
func TestAppleMusicService_GetAuthURL(t *testing.T) {
    service := NewAppleMusicService()

    state := "test-state"
    redirectURL := "http://localhost:3000/callback"

    authURL, err := service.GetAuthURL(state, redirectURL)
    assert.NoError(t, err)
    assert.Contains(t, authURL, "https://appleid.apple.com/auth/authorize")
}

func TestAppleMusicService_SyncUserData(t *testing.T) {
    service := NewAppleMusicService()

    // Mock Apple Music API server
    server := setupMockAppleMusicServer(t)
    defer server.Close()

    tokens := &OAuthTokens{
        AccessToken: "test-access-token",
        TokenType:   "Bearer",
        ExpiresAt:   time.Now().Add(time.Hour),
    }

    result, err := service.SyncUserData(context.Background(), tokens, time.Now().Add(-24*time.Hour))
    assert.NoError(t, err)
    assert.True(t, result.Success)
    assert.Greater(t, len(result.Items), 0)
}
```

### 2. Google Calendar Service (Planned)

```go
// backend/services/calendar/google/service_test.go
func TestGoogleCalendarService_SyncUserData(t *testing.T) {
    service := NewGoogleCalendarService()

    // Mock Google Calendar API server
    server := setupMockGoogleServer(t)
    defer server.Close()

    tokens := &OAuthTokens{
        AccessToken: "test-access-token",
        TokenType:   "Bearer",
        ExpiresAt:   time.Now().Add(time.Hour),
    }

    result, err := service.SyncUserData(context.Background(), tokens, time.Now().Add(-24*time.Hour))
    assert.NoError(t, err)
    assert.True(t, result.Success)
    assert.Greater(t, len(result.Items), 0)

    // Verify calendar events
    for _, item := range result.Items {
        assert.Equal(t, "calendar_event", item.ItemType)
    }
}
```

## Testing Infrastructure

### 1. Test Helpers

```go
// backend/testing/helpers.go
func setupTestDB(t *testing.T) *sqlx.DB {
    db, err := sqlx.Connect("postgres", "postgres://test:test@localhost/syncer_test?sslmode=disable")
    if err != nil {
        t.Skip("Test database not available:", err)
    }

    // Clean up before tests
    db.MustExec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")

    // Run migrations
    err = runMigrations(db, "up")
    if err != nil {
        t.Fatal("Failed to run migrations:", err)
    }

    return db
}

func createTestUser(t *testing.T, db *sqlx.DB, email, passwordHash, fullName string) *User {
    user := &User{
        Email:     email,
        Password:  passwordHash,
        FullName:  fullName,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    _, err := db.NamedExec(`INSERT INTO users (email, password_hash, full_name, created_at, updated_at)
        VALUES (:email, :password_hash, :full_name, :created_at, :updated_at)`, user)
    if err != nil {
        t.Fatal("Failed to create test user:", err)
    }

    err = db.Get(user, "SELECT * FROM users WHERE email = $1", email)
    if err != nil {
        t.Fatal("Failed to retrieve created user:", err)
    }

    return user
}
```

### 2. Mock Services

```go
// backend/testing/mocks.go
type MockService struct {
    name        string
    displayName string
    category    ServiceCategory
    scopes      []string
}

func (m *MockService) Name() string                      { return m.name }
func (m *MockService) DisplayName() string               { return m.displayName }
func (m *MockService) Category() ServiceCategory         { return m.category }
func (m *MockService) RequiredScopes() []string          { return m.scopes }
func (m *MockService) HealthCheck(ctx context.Context) error { return nil }

func (m *MockService) SyncUserData(ctx context.Context, tokens *OAuthTokens, lastSync time.Time) (*SyncResult[map[string]any], error) {
    return &SyncResult[map[string]any]{
        Success: true,
        Items: []SyncItem[map[string]any]{
            {
                ExternalID:   "mock_item_1",
                ItemType:     "mock_item",
                Action:       ActionCreate,
                Data:         map[string]any{"name": "Mock Item 1"},
                LastModified: time.Now(),
            },
        },
        Errors: []SyncError{},
    }, nil
}
```

## CI/CD Integration

### 1. GitHub Actions Workflow

```yaml
# .github/workflows/test.yml
name: Test

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:13
        env:
          POSTGRES_USER: test
          POSTGRES_PASSWORD: test
          POSTGRES_DB: syncer_test
        ports:
          - 5432:5432

    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Run unit tests
        run: go test -v -race -coverprofile=coverage.out ./...

      - name: Run integration tests
        run: go test -v -tags=integration ./...
        env:
          DATABASE_URL: postgres://test:test@localhost:5432/syncer_test?sslmode=disable
```

## Implementation Timeline and Coverage Goals

### Phase 1: Foundation (Week 1-2)

- [ ] Authentication system tests (85% coverage)
- [ ] JWT service tests (95% coverage)
- [ ] Session service tests (90% coverage)
- [ ] Database migration tests (100% coverage)

### Phase 2: Core Services (Week 3-4)

- [ ] Service registry tests (95% coverage)
- [ ] OAuth manager tests (90% coverage)
- [ ] Base service tests (85% coverage)
- [ ] all services tests (80% coverage)

### Phase 3: Sync Engine (Week 5-6)

- [ ] Sync engine tests (90% coverage)
- [ ] Track transformer tests (95% coverage)
- [ ] Cross-service sync tests (85% coverage)
- [ ] API endpoint tests (80% coverage)

### Phase 4: Advanced Features (Week 7-8)

- [ ] Scheduled sync tests (85% coverage)
- [ ] Security tests (90% coverage)

### Phase 5: Production Readiness (Week 9-10)

- [ ] Performance tests (75% coverage)
- [ ] Integration tests (80% coverage)
- [ ] End-to-end tests (60% coverage)
- [ ] CI/CD pipeline setup
- [ ] Coverage reporting

### Coverage Targets

- **Overall Target**: 85% code coverage
- **Critical Paths**: 95% coverage (auth, security, sync engine)
- **API Endpoints**: 90% coverage
- **Business Logic**: 85% coverage
- **Infrastructure**: 75% coverage

This comprehensive testing plan ensures that both implemented and planned features are thoroughly tested, following Go best practices and maintaining the privacy-first principles of the application.
