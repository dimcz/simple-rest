package records_service

import "simple-rest/data"

type ActionRequest struct {
	UserID   int
	RecordID string
}

func (ar *ActionRequest) GetAll() ([]data.Record, error) {
	records, err := data.Records(ar.UserID)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (ar *ActionRequest) DeleteByRecordID() error {
	err := data.DeleteRecordByID(ar.UserID, ar.RecordID)
	return err
}
