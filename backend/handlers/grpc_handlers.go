package handlers

import (
	"context"
	"fmt"
	"log"

	authpb "github.com/clarity/backend/gen/go/auth"
	healthpb "github.com/clarity/backend/gen/go/health"
	aipb "github.com/clarity/backend/gen/go/ai"
	"github.com/clarity/backend/services"
)

// AuthServer implements the gRPC AuthService
type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
	authService *services.AuthService
}

func NewAuthServer(authService *services.AuthService) *AuthServer {
	return &AuthServer{authService: authService}
}

func (as *AuthServer) SendOTP(ctx context.Context, req *authpb.SendOTPRequest) (*authpb.SendOTPResponse, error) {
	otp, err := as.authService.SendOTP(req.Email)
	if err != nil {
		return &authpb.SendOTPResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to send OTP: %v", err),
		}, nil
	}

	// For testing only - in production, remove this
	log.Printf("OTP sent to %s: %s", req.Email, otp)

	return &authpb.SendOTPResponse{
		Success: true,
		Message: "OTP sent to email",
	}, nil
}

func (as *AuthServer) VerifyOTP(ctx context.Context, req *authpb.VerifyOTPRequest) (*authpb.VerifyOTPResponse, error) {
	user, accessToken, refreshToken, err := as.authService.VerifyOTP(req.Email, req.Otp)
	if err != nil {
		return &authpb.VerifyOTPResponse{
			Success: false,
		}, nil
	}

	return &authpb.VerifyOTPResponse{
		Success:      true,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: &authpb.User{
			Id:          user.ID,
			Email:       user.Email,
			Name:        user.Name,
			DateOfBirth: user.DateOfBirth,
			Gender:      user.Gender,
			BloodType:   user.BloodType,
			CreatedAt:   user.CreatedAt.Unix(),
			UpdatedAt:   user.UpdatedAt.Unix(),
		},
	}, nil
}

func (as *AuthServer) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	accessToken, err := as.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &authpb.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken, // In production, rotate this
	}, nil
}

// HealthRecordsServer implements the gRPC HealthRecordsService
type HealthRecordsServer struct {
	healthpb.UnimplementedHealthRecordsServiceServer
	healthService *services.HealthRecordsService
}

func NewHealthRecordsServer(healthService *services.HealthRecordsService) *HealthRecordsServer {
	return &HealthRecordsServer{healthService: healthService}
}

func (hrs *HealthRecordsServer) CreateRecord(ctx context.Context, req *healthpb.CreateRecordRequest) (*healthpb.HealthRecord, error) {
	record, err := hrs.healthService.CreateRecord(req.UserId, req.RecordType, req.Title, req.Description, req.Metadata)
	if err != nil {
		log.Printf("Error creating record: %v", err)
		return nil, err
	}

	return &healthpb.HealthRecord{
		Id:          record.ID,
		UserId:      record.UserID,
		RecordType:  record.RecordType,
		Title:       record.Title,
		Description: record.Description,
		Metadata:    req.Metadata,
		CreatedAt:   record.CreatedAt.String(),
		UpdatedAt:   record.UpdatedAt.String(),
	}, nil
}

func (hrs *HealthRecordsServer) GetRecord(ctx context.Context, req *healthpb.GetRecordRequest) (*healthpb.HealthRecord, error) {
	record, err := hrs.healthService.GetRecord(req.RecordId)
	if err != nil {
		return nil, err
	}

	return &healthpb.HealthRecord{
		Id:          record.ID,
		UserId:      record.UserID,
		RecordType:  record.RecordType,
		Title:       record.Title,
		Description: record.Description,
		CreatedAt:   record.CreatedAt.String(),
		UpdatedAt:   record.UpdatedAt.String(),
	}, nil
}

func (hrs *HealthRecordsServer) ListRecords(ctx context.Context, req *healthpb.ListRecordsRequest) (*healthpb.ListRecordsResponse, error) {
	records, total, err := hrs.healthService.ListRecords(req.UserId, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}

	pbRecords := make([]*healthpb.HealthRecord, len(records))
	for i, record := range records {
		pbRecords[i] = &healthpb.HealthRecord{
			Id:          record.ID,
			UserId:      record.UserID,
			RecordType:  record.RecordType,
			Title:       record.Title,
			Description: record.Description,
			CreatedAt:   record.CreatedAt.String(),
			UpdatedAt:   record.UpdatedAt.String(),
		}
	}

	return &healthpb.ListRecordsResponse{
		Records: pbRecords,
		Total:   int32(total),
	}, nil
}

func (hrs *HealthRecordsServer) UpdateRecord(ctx context.Context, req *healthpb.UpdateRecordRequest) (*healthpb.HealthRecord, error) {
	record, err := hrs.healthService.UpdateRecord(req.RecordId, req.Title, req.Description, req.Metadata)
	if err != nil {
		return nil, err
	}

	return &healthpb.HealthRecord{
		Id:          record.ID,
		UserId:      record.UserID,
		RecordType:  record.RecordType,
		Title:       record.Title,
		Description: record.Description,
		CreatedAt:   record.CreatedAt.String(),
		UpdatedAt:   record.UpdatedAt.String(),
	}, nil
}

func (hrs *HealthRecordsServer) DeleteRecord(ctx context.Context, req *healthpb.DeleteRecordRequest) (*healthpb.DeleteRecordResponse, error) {
	err := hrs.healthService.DeleteRecord(req.RecordId)
	if err != nil {
		return &healthpb.DeleteRecordResponse{Success: false}, nil
	}

	return &healthpb.DeleteRecordResponse{Success: true}, nil
}

// AIServer implements the gRPC AIService
type AIServer struct {
	aipb.UnimplementedAIServiceServer
	aiService *services.AIService
}

func NewAIServer(aiService *services.AIService) *AIServer {
	return &AIServer{aiService: aiService}
}

func (ai *AIServer) ScanPrescription(ctx context.Context, req *aipb.ScanPrescriptionRequest) (*aipb.ScanPrescriptionResponse, error) {
	extractedData, err := ai.aiService.ScanPrescription(req.UserId, req.ImageData)
	if err != nil {
		return &aipb.ScanPrescriptionResponse{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &aipb.ScanPrescriptionResponse{
		Success:       true,
		PrescriptionText: fmt.Sprintf("%v", extractedData),
		ExtractedData: extractedData,
	}, nil
}

func (ai *AIServer) SummarizeHealth(ctx context.Context, req *aipb.SummarizeHealthRequest) (*aipb.SummarizeHealthResponse, error) {
	summary, findings, recommendations, err := ai.aiService.SummarizeHealth(req.UserId, int(req.Days))
	if err != nil {
		return &aipb.SummarizeHealthResponse{
			Success: false,
		}, nil
	}

	return &aipb.SummarizeHealthResponse{
		Success:         true,
		Summary:         summary,
		KeyFindings:     findings,
		Recommendations: recommendations,
	}, nil
}

func (ai *AIServer) DoctorChat(stream aipb.AIService_DoctorChatServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}

		response, err := ai.aiService.DoctorChat(req.UserId, req.ConversationId, req.Message)
		if err != nil {
			log.Printf("Error in doctor chat: %v", err)
			continue
		}

		chatResponse := &aipb.DoctorChatResponse{
			ConversationId: req.ConversationId,
			Response:       response,
			IsAI:           true,
			Timestamp:      int64(0), // Will be set by server
		}

		if err := stream.Send(chatResponse); err != nil {
			return err
		}
	}
}
