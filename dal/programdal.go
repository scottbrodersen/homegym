package dal

import (
	"fmt"

	badger "github.com/dgraph-io/badger/v4"
)

const (
	programKey         = "program"
	programInstanceKey = "programinstance"
	activeProgramKey   = "activeprogram"
)

// AddProgram writes a program to the database.
// When the program already exists it is overwritten.
func (c *DBClient) AddProgram(userID, activityID, programID string, program []byte) error {
	prefix := []string{userKey, userID, activityKey, activityID, programKey, programID}

	programEntry := badger.NewEntry(key(prefix), program)
	updates := []*badger.Entry{programEntry}

	if err := writeUpdates(c, updates); err != nil {
		return fmt.Errorf("failed to add program: %w", err)
	}

	return nil
}

// GetProgramPage gets a page of programs.
// To get a specific program, set pageSize to 1.
// Returns nil when a specific program was requested and not found.
func (c *DBClient) GetProgramPage(userID, activityID, previousProgramID string, pageSize int) ([][]byte, error) {
	prefix := []string{userKey, userID, activityKey, activityID, programKey}
	var startKey []byte = nil
	firstPage := true
	if previousProgramID != "" {
		startKey = key(append(prefix, previousProgramID))
		if pageSize > 1 {
			firstPage = false
		}
	}
	entries, err := readKeyPrefixPage(c, startKey, key(prefix), pageSize, "", firstPage, false)
	if err != nil {
		return nil, fmt.Errorf("failed to read programs: %w", err)
	}

	programs := [][]byte{}
	if pageSize > 1 {

		for _, v := range entries {
			programs = append(programs, v.Value)
		}

		return programs, nil
	}
	// When getting a specific program, make sure we got the right one
	if len(entries) == 0 {
		return nil, nil
	}

	if string(entries[0].Key) != string(startKey) {
		return nil, nil
	}
	programs = append(programs, entries[0].Value)

	return programs, nil
}

// AddProgramInstance adds or updates a program instance.
// Sets the instance as active when activityID is not an empty string.
func (c *DBClient) AddProgramInstance(userID, programID, instanceID, activityID string, instance []byte) error {
	instancePrefix := []string{userKey, userID, programKey, programID, programInstanceKey, instanceID}

	instanceEntry := badger.NewEntry(key(instancePrefix), instance)

	updates := []*badger.Entry{instanceEntry}

	// if activityID != "" {

	// 	activePrefix := []string{userKey, userID, activityKey, activityID, activeProgramKey}

	// 	// value is in format {instanceID}:{programID}
	// 	activeProgramEntry := badger.NewEntry(key(activePrefix), []byte(fmt.Sprintf("%s:%s", instanceID, programID)))

	// 	updates = append(updates, activeProgramEntry)
	// }

	if err := writeUpdates(c, updates); err != nil {
		return fmt.Errorf("failed to add program instance: %w", err)
	}

	return nil
}

// GetProgramInstancePage gets a page of program instances.
// To get a specific instance, set pageSize to 1 and previousProgramInstanceID to the instance ID.
// Returns nil if a specific instance was requested and not found.
func (c *DBClient) GetProgramInstancePage(userID, programID, previousProgramInstanceID string, pageSize int) ([][]byte, error) {
	prefix := []string{userKey, userID, programKey, programID, programInstanceKey}
	var startKey []byte = nil
	firstPage := true
	if previousProgramInstanceID != "" {
		startKey = key(append(prefix, previousProgramInstanceID))
		if pageSize > 1 {
			firstPage = false
		}
	}
	entries, err := readKeyPrefixPage(c, startKey, key(prefix), pageSize, "", firstPage, false)
	if err != nil {
		return nil, fmt.Errorf("failed to read program instances: %w", err)
	}
	programInstances := [][]byte{}

	if pageSize > 1 {
		for _, v := range entries {
			programInstances = append(programInstances, v.Value)
		}

		return programInstances, nil
	}

	// When getting a specific instance, make sure we got the right one
	if len(entries) == 0 {
		return nil, nil
	}

	if string(entries[0].Key) != string(startKey) {
		return nil, nil
	}
	programInstances = append(programInstances, entries[0].Value)

	return programInstances, nil
}

// func (c *DBClient) SetActiveProgramInstance(userID, activityID, programID, instanceID string) error {
// 	prefix := []string{userKey, userID, activityKey, activityID, activeProgramKey}

// 	// value is in format {instanceID}:{programID}
// 	activeProgramEntry := badger.NewEntry(key(prefix), []byte(fmt.Sprintf("%s:%s", instanceID, programID)))
// 	updates := []*badger.Entry{activeProgramEntry}

// 	if err := writeUpdates(c, updates); err != nil {
// 		return fmt.Errorf("failed to set active program: %w", err)
// 	}

// 	return nil
// }

func (c *DBClient) ActivateProgramInstance(userID, activityID, programID, instanceID string) error {
	prefix := []string{userKey, userID, activityKey, activityID, activeProgramKey, instanceID}

	// value is {programID:instanceID}
	activeProgramEntry := badger.NewEntry(key(prefix), []byte(fmt.Sprintf("%s:%s", programID, instanceID)))
	updates := []*badger.Entry{activeProgramEntry}

	if err := writeUpdates(c, updates); err != nil {
		return fmt.Errorf("failed to activate program: %w", err)
	}

	return nil
}

// GetActiveProgramInstancePage gets a page of active program instance IDs with their programID for a specific activity.
// Returned values are an array of byte arrays of value {programID}:{programInstanceID}.
// To get a specific instance, set pageSize to 1.
// Returns nil if a specific instance was requested and not found.
func (c *DBClient) GetActiveProgramInstancePage(userID, activityID, previousActiveInstanceID string, pageSize int) ([][]byte, error) {
	prefix := []string{userKey, userID, activityKey, activityID, activeProgramKey}
	var startKey []byte = nil
	firstPage := true
	if previousActiveInstanceID != "" {
		startKey = key(append(prefix, previousActiveInstanceID))
		if pageSize > 1 {
			firstPage = false
		}
	}
	entries, err := readKeyPrefixPage(c, startKey, key(prefix), pageSize, "", firstPage, false)
	if err != nil {
		return nil, fmt.Errorf("failed to read program instances: %w", err)
	}
	activeInstances := [][]byte{}

	if pageSize > 1 {
		for _, v := range entries {
			activeInstances = append(activeInstances, v.Value)
		}

		return activeInstances, nil
	}

	// When getting a specific instance, make sure we got the right one
	if len(entries) == 0 {
		return nil, nil
	}

	if string(entries[0].Key) != string(startKey) {
		return nil, nil
	}
	activeInstances = append(activeInstances, entries[0].Value)

	return activeInstances, nil
}

// func (c *DBClient) GetActiveProgramInstance(userID, activityID string) ([]byte, error) {

// 	prefix := []string{userKey, userID, activityKey, activityID, activeProgramKey}
// 	entry, err := readItem(c, key(prefix))

// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read event: %w", err)
// 	}

// 	if entry == nil {
// 		return nil, nil
// 	}

// 	ids := strings.Split(string(entry.Value), ":")

// 	instancePage, err := c.GetProgramInstancePage(userID, ids[1], ids[0], 1)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read program instance: %w", err)
// 	}

// 	if len(instancePage) == 0 {
// 		return nil, nil
// 	}

// 	return instancePage[0], nil

// }

func (c *DBClient) DeactivateProgramInstance(userID, activityID, activeInstanceID string) error {

	keyPrefix := []string{userKey, userID, activityKey, activityID, activeProgramKey, activeInstanceID}
	keys := [][]byte{key(keyPrefix)}
	err := deleteItems(c, keys)
	if err != nil {
		return fmt.Errorf("failed to delete active program: %w", err)
	}
	return nil
}
