package handler

import (
	"bytes"
	"context" // 👈 1. 新增
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time" // 👈 2. 新增 (為了初始化 token payload)

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kdotwei/hpl-scoreboard/internal/db"
	"github.com/kdotwei/hpl-scoreboard/internal/middleware" // 👈 3. 新增
	"github.com/kdotwei/hpl-scoreboard/internal/service"
	"github.com/kdotwei/hpl-scoreboard/internal/service/mocks"
	"github.com/kdotwei/hpl-scoreboard/internal/token" // 👈 4. 新增
	token_mocks "github.com/kdotwei/hpl-scoreboard/internal/token/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateScore(t *testing.T) {
	testCases := []struct {
		name              string
		requestBody       CreateScoreRequest
		mockUser          string
		expectedStatus    int
		shouldCallService bool
		setupMock         func(*mocks.Service, CreateScoreRequest, string)
	}{
		{
			name: "successful score creation with all fields",
			requestBody: CreateScoreRequest{
				Gflops:        1234.56,
				ProblemSizeN:  20000,
				BlockSizeNb:   512,
				LinuxUsername: "hpl_user",
				N:             20000,
				NB:            512,
				P:             4,
				Q:             4,
				ExecutionTime: 125.75,
			},
			mockUser:          "jwt-user",
			expectedStatus:    http.StatusCreated,
			shouldCallService: true,
			setupMock: func(mockService *mocks.Service, req CreateScoreRequest, user string) {
				mockService.On("CreateScore", mock.Anything, mock.MatchedBy(func(arg service.CreateScoreParams) bool {
					return arg.UserID == user &&
						arg.Gflops == req.Gflops &&
						arg.ProblemSizeN == req.ProblemSizeN &&
						arg.BlockSizeNb == req.BlockSizeNb &&
						arg.LinuxUsername == req.LinuxUsername &&
						arg.N == req.N &&
						arg.NB == req.NB &&
						arg.P == req.P &&
						arg.Q == req.Q &&
						arg.ExecutionTime == req.ExecutionTime
				})).Return(&db.Score{
					ID:            pgtype.UUID{Bytes: [16]byte{1, 2, 3}, Valid: true},
					UserID:        user,
					Gflops:        req.Gflops,
					ProblemSizeN:  int32(req.ProblemSizeN),
					BlockSizeNb:   int32(req.BlockSizeNb),
					LinuxUsername: req.LinuxUsername,
					N:             int32(req.N),
					Nb:            int32(req.NB),
					P:             int32(req.P),
					Q:             int32(req.Q),
					ExecutionTime: req.ExecutionTime,
					SubmittedAt:   time.Now(),
				}, nil)
			},
		},
		{
			name: "minimal valid score creation",
			requestBody: CreateScoreRequest{
				Gflops:        100.0,
				ProblemSizeN:  1000,
				BlockSizeNb:   64,
				LinuxUsername: "minimal_user",
				N:             1000,
				NB:            64,
				P:             1,
				Q:             1,
				ExecutionTime: 50.0,
			},
			mockUser:          "minimal-jwt-user",
			expectedStatus:    http.StatusCreated,
			shouldCallService: true,
			setupMock: func(mockService *mocks.Service, req CreateScoreRequest, user string) {
				mockService.On("CreateScore", mock.Anything, mock.Anything).Return(&db.Score{
					ID:            pgtype.UUID{Bytes: [16]byte{4, 5, 6}, Valid: true},
					UserID:        user,
					Gflops:        req.Gflops,
					LinuxUsername: req.LinuxUsername,
					SubmittedAt:   time.Now(),
				}, nil)
			},
		},
		{
			name: "high performance score with large matrices",
			requestBody: CreateScoreRequest{
				Gflops:        9876.54,
				ProblemSizeN:  100000,
				BlockSizeNb:   1024,
				LinuxUsername: "supercomputer_user",
				N:             100000,
				NB:            1024,
				P:             8,
				Q:             8,
				ExecutionTime: 3600.25,
			},
			mockUser:          "performance-user",
			expectedStatus:    http.StatusCreated,
			shouldCallService: true,
			setupMock: func(mockService *mocks.Service, req CreateScoreRequest, user string) {
				mockService.On("CreateScore", mock.Anything, mock.Anything).Return(&db.Score{
					ID:            pgtype.UUID{Bytes: [16]byte{7, 8, 9}, Valid: true},
					UserID:        user,
					Gflops:        req.Gflops,
					ExecutionTime: req.ExecutionTime,
					SubmittedAt:   time.Now(),
				}, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 1. Setup Mock
			mockService := new(mocks.Service)
			mockTokenMaker := new(token_mocks.Maker)
			h := NewHandler(mockService, mockTokenMaker)

			// Setup service mock expectations
			if tc.shouldCallService {
				tc.setupMock(mockService, tc.requestBody, tc.mockUser)
			}

			// 2. Prepare request body
			jsonBody, err := json.Marshal(tc.requestBody)
			assert.NoError(t, err)

			// 3. Create HTTP Request
			req, err := http.NewRequest("POST", "/api/v1/scores", bytes.NewBuffer(jsonBody))
			assert.NoError(t, err)

			// 4. Inject Auth Payload into Context (simulating middleware)
			mockPayload := &token.Payload{
				Username:  tc.mockUser,
				IssuedAt:  time.Now(),
				ExpiredAt: time.Now().Add(time.Hour),
			}
			ctx := context.WithValue(req.Context(), middleware.AuthorizationPayloadKey, mockPayload)
			req = req.WithContext(ctx)

			// 5. Execute Handler
			rr := httptest.NewRecorder()
			http.HandlerFunc(h.CreateScore).ServeHTTP(rr, req)

			// 6. Assertions
			assert.Equal(t, tc.expectedStatus, rr.Code)

			if tc.shouldCallService && tc.expectedStatus == http.StatusCreated {
				// Verify response contains the created score
				var response db.Score
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.mockUser, response.UserID)
				assert.Equal(t, tc.requestBody.Gflops, response.Gflops)
			}

			// Verify all mock expectations were met
			mockService.AssertExpectations(t)
		})
	}
}

func TestCreateScore_ErrorCases(t *testing.T) {
	testCases := []struct {
		name           string
		requestBody    string // Using string to test malformed JSON
		mockUser       string
		hasAuthPayload bool
		expectedStatus int
		setupMock      func(*mocks.Service)
	}{
		{
			name:           "invalid JSON body",
			requestBody:    `{"gflops": "invalid", "problem_size_n": 1000}`,
			mockUser:       "test-user",
			hasAuthPayload: true,
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(mockService *mocks.Service) {},
		},
		{
			name:           "missing authorization payload",
			requestBody:    `{"gflops": 123.45, "problem_size_n": 1000, "block_size_nb": 256, "linux_username": "test", "n": 1000, "nb": 256, "p": 1, "q": 1, "execution_time": 50.0}`,
			mockUser:       "",
			hasAuthPayload: false,
			expectedStatus: http.StatusUnauthorized,
			setupMock:      func(mockService *mocks.Service) {},
		},
		{
			name:           "service layer error",
			requestBody:    `{"gflops": 123.45, "problem_size_n": 1000, "block_size_nb": 256, "linux_username": "test", "n": 1000, "nb": 256, "p": 1, "q": 1, "execution_time": 50.0}`,
			mockUser:       "test-user",
			hasAuthPayload: true,
			expectedStatus: http.StatusInternalServerError,
			setupMock: func(mockService *mocks.Service) {
				mockService.On("CreateScore", mock.Anything, mock.Anything).Return(nil, assert.AnError)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 1. Setup Mock
			mockService := new(mocks.Service)
			mockTokenMaker := new(token_mocks.Maker)
			h := NewHandler(mockService, mockTokenMaker)

			// Setup service mock expectations
			tc.setupMock(mockService)

			// 2. Create HTTP Request
			req, err := http.NewRequest("POST", "/api/v1/scores", bytes.NewBufferString(tc.requestBody))
			assert.NoError(t, err)

			// 3. Conditionally inject Auth Payload into Context
			if tc.hasAuthPayload {
				mockPayload := &token.Payload{
					Username:  tc.mockUser,
					IssuedAt:  time.Now(),
					ExpiredAt: time.Now().Add(time.Hour),
				}
				ctx := context.WithValue(req.Context(), middleware.AuthorizationPayloadKey, mockPayload)
				req = req.WithContext(ctx)
			}

			// 4. Execute Handler
			rr := httptest.NewRecorder()
			http.HandlerFunc(h.CreateScore).ServeHTTP(rr, req)

			// 5. Assertions
			assert.Equal(t, tc.expectedStatus, rr.Code)

			// Verify all mock expectations were met
			mockService.AssertExpectations(t)
		})
	}
}

func TestListScores(t *testing.T) {
	testCases := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedLimit  int32
		expectedOffset int32
		setupMock      func(*mocks.Service)
		expectedScores []db.Score
	}{
		{
			name:           "successful list scores with default limit and offset",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			expectedLimit:  10,
			expectedOffset: 0,
			setupMock: func(mockService *mocks.Service) {
				mockScores := []db.Score{
					{
						ID:            pgtype.UUID{Bytes: [16]byte{1, 2, 3}, Valid: true},
						UserID:        "user1",
						Gflops:        1000.0,
						ProblemSizeN:  50000,
						BlockSizeNb:   256,
						LinuxUsername: "hpl_user1",
						N:             50000,
						Nb:            256,
						P:             4,
						Q:             4,
						ExecutionTime: 1800.5,
						SubmittedAt:   time.Now(),
					},
					{
						ID:            pgtype.UUID{Bytes: [16]byte{4, 5, 6}, Valid: true},
						UserID:        "user2",
						Gflops:        800.0,
						ProblemSizeN:  40000,
						BlockSizeNb:   128,
						LinuxUsername: "hpl_user2",
						N:             40000,
						Nb:            128,
						P:             2,
						Q:             2,
						ExecutionTime: 1200.3,
						SubmittedAt:   time.Now(),
					},
				}
				mockService.On("ListScores", mock.Anything, int32(10), int32(0)).Return(mockScores, nil)
			},
		},
		{
			name:           "successful list scores with custom limit",
			queryParams:    "?limit=5",
			expectedStatus: http.StatusOK,
			expectedLimit:  5,
			expectedOffset: 0,
			setupMock: func(mockService *mocks.Service) {
				mockScores := []db.Score{
					{
						ID:            pgtype.UUID{Bytes: [16]byte{7, 8, 9}, Valid: true},
						UserID:        "user3",
						Gflops:        1200.0,
						ProblemSizeN:  60000,
						BlockSizeNb:   512,
						LinuxUsername: "hpl_user3",
						N:             60000,
						Nb:            512,
						P:             8,
						Q:             8,
						ExecutionTime: 2400.7,
						SubmittedAt:   time.Now(),
					},
				}
				mockService.On("ListScores", mock.Anything, int32(5), int32(0)).Return(mockScores, nil)
			},
		},
		{
			name:           "successful list scores with offset",
			queryParams:    "?limit=5&offset=10",
			expectedStatus: http.StatusOK,
			expectedLimit:  5,
			expectedOffset: 10,
			setupMock: func(mockService *mocks.Service) {
				mockScores := []db.Score{
					{
						ID:            pgtype.UUID{Bytes: [16]byte{11, 12, 13}, Valid: true},
						UserID:        "user4",
						Gflops:        500.0,
						ProblemSizeN:  30000,
						BlockSizeNb:   128,
						LinuxUsername: "hpl_user4",
						N:             30000,
						Nb:            128,
						P:             2,
						Q:             2,
						ExecutionTime: 600.3,
						SubmittedAt:   time.Now(),
					},
				}
				mockService.On("ListScores", mock.Anything, int32(5), int32(10)).Return(mockScores, nil)
			},
		},
		{
			name:           "invalid limit returns bad request",
			queryParams:    "?limit=invalid",
			expectedStatus: http.StatusBadRequest,
			expectedLimit:  0,
			expectedOffset: 0,
			setupMock: func(mockService *mocks.Service) {
				// No mock call expected for bad request
			},
		},
		{
			name:           "negative limit returns bad request",
			queryParams:    "?limit=-1",
			expectedStatus: http.StatusBadRequest,
			expectedLimit:  0,
			expectedOffset: 0,
			setupMock: func(mockService *mocks.Service) {
				// No mock call expected for bad request
			},
		},
		{
			name:           "invalid offset returns bad request",
			queryParams:    "?offset=invalid",
			expectedStatus: http.StatusBadRequest,
			expectedLimit:  0,
			expectedOffset: 0,
			setupMock: func(mockService *mocks.Service) {
				// No mock call expected for bad request
			},
		},
		{
			name:           "negative offset returns bad request",
			queryParams:    "?offset=-1",
			expectedStatus: http.StatusBadRequest,
			expectedLimit:  0,
			expectedOffset: 0,
			setupMock: func(mockService *mocks.Service) {
				// No mock call expected for bad request
			},
		},
		{
			name:           "service layer error",
			queryParams:    "",
			expectedStatus: http.StatusInternalServerError,
			expectedLimit:  10,
			expectedOffset: 0,
			setupMock: func(mockService *mocks.Service) {
				mockService.On("ListScores", mock.Anything, int32(10), int32(0)).Return(nil, assert.AnError)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 1. Setup Mock
			mockService := new(mocks.Service)
			mockTokenMaker := new(token_mocks.Maker)
			h := NewHandler(mockService, mockTokenMaker)

			// Setup service mock expectations
			tc.setupMock(mockService)

			// 2. Create HTTP Request (no authentication required)
			req, err := http.NewRequest("GET", "/api/v1/scores"+tc.queryParams, nil)
			assert.NoError(t, err)

			// 3. Execute Handler
			rr := httptest.NewRecorder()
			http.HandlerFunc(h.ListScores).ServeHTTP(rr, req)

			// 4. Assertions
			assert.Equal(t, tc.expectedStatus, rr.Code)

			if tc.expectedStatus == http.StatusOK {
				var scores []db.Score
				err := json.Unmarshal(rr.Body.Bytes(), &scores)
				assert.NoError(t, err)
			}

			// Verify all mock expectations were met
			mockService.AssertExpectations(t)
		})
	}
}

func TestListScoresWithPagination(t *testing.T) {
	testCases := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedLimit  int32
		expectedOffset int32
		setupMock      func(*mocks.Service)
	}{
		{
			name:           "successful pagination with default parameters",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			expectedLimit:  10,
			expectedOffset: 0,
			setupMock: func(mockService *mocks.Service) {
				mockResponse := &service.PaginatedScoresResponse{
					Scores: []db.Score{
						{ID: pgtype.UUID{Valid: true}, Gflops: 100.0, UserID: "user1"},
					},
					HasMore:      false,
					TotalRecords: 1,
					Limit:        10,
					Offset:       0,
				}
				mockService.On("ListScoresWithPagination", mock.Anything, mock.MatchedBy(func(params service.ListScoresParams) bool {
					return params.Limit == 10 && params.Offset == 0
				})).Return(mockResponse, nil)
			},
		},
		{
			name:           "successful pagination with custom limit",
			queryParams:    "?limit=5",
			expectedStatus: http.StatusOK,
			expectedLimit:  5,
			expectedOffset: 0,
			setupMock: func(mockService *mocks.Service) {
				mockResponse := &service.PaginatedScoresResponse{
					Scores: []db.Score{
						{ID: pgtype.UUID{Valid: true}, Gflops: 200.0, UserID: "user2"},
					},
					HasMore:      true,
					TotalRecords: 50,
					Limit:        5,
					Offset:       0,
				}
				mockService.On("ListScoresWithPagination", mock.Anything, mock.MatchedBy(func(params service.ListScoresParams) bool {
					return params.Limit == 5 && params.Offset == 0
				})).Return(mockResponse, nil)
			},
		},
		{
			name:           "successful pagination with limit and offset",
			queryParams:    "?limit=5&offset=10",
			expectedStatus: http.StatusOK,
			expectedLimit:  5,
			expectedOffset: 10,
			setupMock: func(mockService *mocks.Service) {
				mockResponse := &service.PaginatedScoresResponse{
					Scores: []db.Score{
						{ID: pgtype.UUID{Valid: true}, Gflops: 300.0, UserID: "user3"},
					},
					HasMore:      true,
					TotalRecords: 50,
					Limit:        5,
					Offset:       10,
				}
				mockService.On("ListScoresWithPagination", mock.Anything, mock.MatchedBy(func(params service.ListScoresParams) bool {
					return params.Limit == 5 && params.Offset == 10
				})).Return(mockResponse, nil)
			},
		},
		{
			name:           "invalid limit returns bad request",
			queryParams:    "?limit=invalid",
			expectedStatus: http.StatusBadRequest,
			expectedLimit:  0,
			expectedOffset: 0,
			setupMock: func(mockService *mocks.Service) {
				// No mock call expected for bad request
			},
		},
		{
			name:           "invalid offset returns bad request",
			queryParams:    "?offset=invalid",
			expectedStatus: http.StatusBadRequest,
			expectedLimit:  0,
			expectedOffset: 0,
			setupMock: func(mockService *mocks.Service) {
				// No mock call expected for bad request
			},
		},
		{
			name:           "service error",
			queryParams:    "?limit=5",
			expectedStatus: http.StatusInternalServerError,
			expectedLimit:  5,
			expectedOffset: 0,
			setupMock: func(mockService *mocks.Service) {
				mockService.On("ListScoresWithPagination", mock.Anything, mock.AnythingOfType("service.ListScoresParams")).Return(nil, assert.AnError)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 1. Setup Mock
			mockService := new(mocks.Service)
			mockTokenMaker := new(token_mocks.Maker)
			h := NewHandler(mockService, mockTokenMaker)

			// Setup service mock expectations
			tc.setupMock(mockService)

			// 2. Create HTTP Request
			req, err := http.NewRequest("GET", "/api/v1/scores/paginated"+tc.queryParams, nil)
			assert.NoError(t, err)

			// 3. Execute Handler
			rr := httptest.NewRecorder()
			http.HandlerFunc(h.ListScoresWithPagination).ServeHTTP(rr, req)

			// 4. Assertions
			assert.Equal(t, tc.expectedStatus, rr.Code)

			if tc.expectedStatus == http.StatusOK {
				var response service.PaginatedScoresResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotNil(t, response.Scores)
			}

			// Verify all mock expectations were met
			mockService.AssertExpectations(t)
		})
	}
}
