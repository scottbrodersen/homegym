package dal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	badger "github.com/dgraph-io/badger/v4"
)

type Dstore interface {
	AddActivity(userID, activityID, activityName string) error
	ReadActivity(userID, activityID string) (*string, []string, error)
	GetActivityNames(userID string) (map[string]string, error)
	UpdateActivity(userID, activityID, activityName string) error
	AddExerciseToActivity(userID, activityID, exerciseID string) error
	UpdateActivityExercises(userID, activityID string, exIDsToAdd, exIDsToDelete []string) error

	NewUser(id, email, pwdHash, pwdHashVersion, role string) error
	ReadUser(id string) (*string, *string, *string, *string, error)
	UpdateUserProfile(id, email string) error
	UpdateUserPassword(id, pwdHash, pwdVersion string) error
	ChangeUserRole(id, role string) error
	UpdatePwdVersion(userID, version string) error

	AddExercise(userID, exerciseID string, exercise []byte) error
	UpdateExercise(userID, exerciseID string, exercise []byte) error
	GetExercise(userID, exerciseID string) ([]byte, error)
	GetExercises(userID string) ([][]byte, error)

	AddEvent(userID, eventID, activityID string, date int64, event []byte) error
	ShiftEvent(userID, eventID, activityID string, currentDate, newDate int64, event []byte) error
	GetEvent(userID, eventID string, eventDate int64) ([]byte, error)
	GetEventPage(userID, previousEventID string, previousDate int64, pageSize int) ([][]byte, error)
	GetEventExercises(userID, eventID string) ([][]byte, error)
	AddExercisesToEvent(userID, eventID string, exerciseIDs map[int]string, exerciseInstances map[int][]byte) error

	AddProgram(userID, activityID, programID string, program []byte) error
	GetProgramPage(userID, activityID, previousProgramID string, pageSize int) ([][]byte, error)
	AddProgramInstance(userID, activityID, programID, instanceID string, instance []byte) error
	GetProgramInstancePage(userID, activityID, programID, instanceID string, pageSize int) ([][]byte, error)
	SetActiveProgramInstance(userID, activityID, programID, instanceID string) error
	GetActiveProgramInstance(userID, activityID, programID string) ([]byte, error)

	Destroy()
	GetKeys(usage string) (map[string][]byte, map[string][]byte, error)
	RotateKeys(newRSAKey []byte, keyID, usage string) error
	DeleteKey(keyID, usage string) error

	AddSession(username, sessionID string, ttl int) error
	GetSession(sessionID string) (*string, *int64, error)
	DeleteSession(sessionID string) error
	GetSessionExpiries() (map[string]int64, error)

	// Iter8er()
}

var (
	ErrNotUnique error = errors.New("value not unique")
)

const (
	userKey           = "user"
	idKey             = "id"
	roleKey           = "role"
	nameKey           = "name"
	emailKey          = "email"
	passHashKey       = "phash"
	versionKey        = "version"
	activityKey       = "activity"
	exerciseKey       = "exercise"
	typeKey           = "type"
	indexKey          = "index"
	eventKey          = "event"
	instanceKey       = "instance"
	sessionKey        = "session"
	expiresKey        = "expires"
	tokenCryptoKeyKey = "tokenkey"
)

type DBClient struct {
	db   *badger.DB
	path string
}

var DB Dstore

// Opens the database, creating it if necessary, and returns a client.
func InitClient(path string) (*DBClient, error) {
	dalClient := &DBClient{path: path}
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	dalClient.db = db
	DB = dalClient
	return dalClient, nil
}

func (c *DBClient) Destroy() {
	log.Info("closing database")
	c.db.Close()
}

// Adds items that represent a user
func (c *DBClient) NewUser(id, email, pwdHash, pwdHashVersion, role string) error {
	prefix := []string{userKey, id}
	idKey := key(append(prefix, "id"))

	// Check if user ID is already used
	existingID, err := readItem(c, idKey)
	if err != nil {
		return fmt.Errorf("failed to get user item: %w", err)
	}
	if existingID != nil {
		return ErrNotUnique
	}
	idEntry := badger.NewEntry(idKey, []byte(id))
	roleEntry := badger.NewEntry(key(append(prefix, roleKey)), []byte(role))
	emailEntry := badger.NewEntry(key(append(prefix, emailKey)), []byte(email))
	hashEntry := badger.NewEntry(key(append(prefix, passHashKey)), []byte(pwdHash))
	versionEntry := badger.NewEntry(key(append(prefix, versionKey)), []byte(pwdHashVersion))

	if err := writeUpdates(c, []*badger.Entry{idEntry, roleEntry, emailEntry, hashEntry, versionEntry}); err != nil {
		return fmt.Errorf("failed to add user: %w", err)
	}
	return nil
}

func (c *DBClient) ReadUser(id string) (*string, *string, *string, *string, error) {
	// return values
	var email, pwdHash, pwdHashVersion, role string
	props, err := readKeyPrefix(c, keyPrefix([]string{userKey, id}))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error getting user: %w", err)
	}
	if props == nil {
		// user not found
		return nil, nil, nil, nil, nil
	}
	for _, v := range props {
		itemKey := strings.Split(string(v.Key), "#")[strings.Count(string(v.Key), "#")]
		switch itemKey {

		case roleKey:
			role = string(v.Value)
		case emailKey:
			email = string(v.Value)
		case passHashKey:
			pwdHash = string(v.Value)
		case versionKey:
			pwdHashVersion = string(v.Value)
		}
	}

	return &email, &pwdHash, &pwdHashVersion, &role, nil
}

func (c *DBClient) UpdateUserProfile(id, email string) error {
	prefix := []string{userKey, id}

	emailEntry := badger.NewEntry(key(append(prefix, emailKey)), []byte(email))

	if err := writeUpdates(c, []*badger.Entry{emailEntry}); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (c *DBClient) UpdateUserPassword(id, pwdHash, pwdVersion string) error {
	prefix := []string{userKey, id}

	pwdVersionEntry := badger.NewEntry(key(append(prefix, versionKey)), []byte(pwdVersion))
	hashEntry := badger.NewEntry(key(append(prefix, passHashKey)), []byte(pwdHash))

	if err := writeUpdates(c, []*badger.Entry{pwdVersionEntry, hashEntry}); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (c *DBClient) ChangeUserRole(id, role string) error {
	prefix := []string{userKey, id}

	roleEntry := badger.NewEntry(key(append(prefix, roleKey)), []byte(role))

	if err := writeUpdates(c, []*badger.Entry{roleEntry}); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// TODO: restrict to admin
func (c *DBClient) UpdatePwdVersion(userID, version string) error {
	prefix := []string{userKey, userID}

	versionEntry := badger.NewEntry(key(append(prefix, versionKey)), []byte(version))

	if err := writeUpdates(c, []*badger.Entry{versionEntry}); err != nil {
		return fmt.Errorf("failed to update user hash version: %w", err)
	}

	return nil
}

// Adds items that represent a user's activity.
// Use [AddExercise] to add exercises to the activity.
func (c *DBClient) AddActivity(userID, activityID, activityName string) error {
	prefix := []string{userKey, userID, activityKey, activityID, nameKey}

	nameEntry := badger.NewEntry(key(prefix), []byte(activityName))
	updates := []*badger.Entry{nameEntry}

	if err := writeUpdates(c, updates); err != nil {
		return fmt.Errorf("failed to add activity: %w", err)
	}

	return nil
}

func (c *DBClient) UpdateActivity(userID, activityID, activityName string) error {
	activityEntry := badger.NewEntry(key([]string{userKey, userID, activityKey, activityID, nameKey}), []byte(activityName))
	if err := writeUpdates(c, []*badger.Entry{activityEntry}); err != nil {
		return fmt.Errorf("failed to update activity: %w", err)
	}

	return nil
}

// Returns a map with keys activity IDs and values activity names.
func (c *DBClient) GetActivityNames(userID string) (map[string]string, error) {
	prefix := []string{userKey, userID, activityKey}
	activityEntries, err := readKeyPrefix(c, keyPrefix(prefix))
	if err != nil {
		return nil, fmt.Errorf("failed to read activities: %w", err)
	}
	var activities map[string]string = map[string]string{}
	rx := regexp.MustCompile("^user:[^#]+#activity:([^#]+)#name$")

	for _, v := range activityEntries {
		if aIDParts := rx.FindStringSubmatch(string(v.Key)); aIDParts != nil {
			activities[aIDParts[1]] = string(v.Value)
			continue
		}
	}

	return activities, nil
}

// Returns the activity name and a slice of exercise IDs
func (c *DBClient) ReadActivity(userID, activityID string) (*string, []string, error) {
	prefix := []string{userKey, userID, activityKey, activityID}
	activityEntries, err := readKeyPrefix(c, keyPrefix(prefix))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read activity: %w", err)
	}
	if len(activityEntries) == 0 {
		return nil, nil, nil
	}
	rAct := regexp.MustCompile(fmt.Sprintf("^%s:[^#]+#%s:[^#]+#%s$", userKey, activityKey, nameKey))
	rEx := regexp.MustCompile(fmt.Sprintf("^%s:[^#]+#%s:[^#]+#%s:.*$", userKey, activityKey, exerciseKey))

	var name string
	var exerciseIDs []string
	for _, v := range activityEntries {
		if rAct.Match(v.Key) {
			name = string(v.Value)
			continue
		}
		if exIDPart := rEx.FindStringSubmatch(string(v.Key)); exIDPart != nil {
			exerciseIDs = append(exerciseIDs, string(v.Value))
			continue
		}
	}
	return &name, exerciseIDs, nil
}

func (c *DBClient) AddExerciseToActivity(userID, activityID, exerciseID string) error {
	prefix := []string{userKey, userID, activityKey, activityID, exerciseKey, exerciseID}

	exID := badger.NewEntry(key(prefix), []byte(exerciseID))
	updates := []*badger.Entry{exID}

	if err := writeUpdates(c, updates); err != nil {
		return fmt.Errorf("failed to add exercise to activity: %w", err)
	}

	return nil
}

func (c *DBClient) UpdateActivityExercises(userID, activityID string, exIDsToAdd, exIDsToDelete []string) error {
	prefix := []string{userKey, userID, activityKey, activityID, exerciseKey}
	updates := []*badger.Entry{}
	deletes := [][]byte{}
	for _, exID := range exIDsToAdd {
		exID := badger.NewEntry(key(append(prefix, exID)), []byte(exID))
		updates = append(updates, exID)
	}

	for _, exID := range exIDsToDelete {
		deletes = append(deletes, key(append(prefix, exID)))
	}

	err := updateDeleteItems(c, updates, deletes)
	if err != nil {
		return fmt.Errorf("failed to update activity exercises: %w", err)
	}

	return nil
}

func (c *DBClient) AddExercise(userID, exerciseID string, exercise []byte) error {
	prefix := []string{userKey, userID, exerciseKey, exerciseID, typeKey}

	ex := badger.NewEntry(key(prefix), exercise)

	updates := []*badger.Entry{ex}

	if err := writeUpdates(c, updates); err != nil {
		return fmt.Errorf("failed to add exercise: %w", err)
	}
	return nil
}

// TODO: remove?
func (c *DBClient) UpdateExercise(userID, exerciseID string, exercise []byte) error {
	return c.AddExercise(userID, exerciseID, exercise)
}

// GetExercise returns an exercise type
func (c *DBClient) GetExercise(userID, exerciseID string) ([]byte, error) {
	prefix := []string{userKey, userID, exerciseKey, exerciseID, typeKey}
	exerciseEntry, err := readItem(c, key(prefix))
	if err != nil {
		return nil, fmt.Errorf("failed to read exercise: %w", err)
	}

	if exerciseEntry == nil {
		return nil, nil
	}
	return exerciseEntry.Value, nil
}

// GetExercises returns a slice of exercise types
func (c *DBClient) GetExercises(userID string) ([][]byte, error) {
	prefix := []string{userKey, userID, exerciseKey}
	exerciseEntries, err := readKeyPrefix(c, keyPrefix(prefix))
	if err != nil {
		return nil, fmt.Errorf("failed to read exercise: %w", err)
	}

	if len(exerciseEntries) == 0 {
		return [][]byte{}, nil
	}

	types := [][]byte{}
	rx := regexp.MustCompile(fmt.Sprintf("^%s:[^#]+#%s:[^#]+#%s$", userKey, exerciseKey, typeKey))

	for _, v := range exerciseEntries {
		if rx.MatchString(string(v.Key)) {
			types = append(types, v.Value)
		}
	}

	return types, nil
}

func (c *DBClient) AddEvent(userID, eventID, activityID string, date int64, event []byte) error {
	prefix := []string{userKey, userID, eventKey, fmt.Sprint(date), idKey, eventID, activityKey, activityID}
	// user:{id}#event:{date}#id:{id}#activity:{activityID}
	eventEntry := badger.NewEntry(key(prefix), event)

	updates := []*badger.Entry{eventEntry}

	if err := writeUpdates(c, updates); err != nil {
		return fmt.Errorf("failed to add event: %w", err)
	}
	return nil
}

// ShiftEvent updates an event that has a new date.
// The date is part of the key so the existing must be deleted and the updated event is created.
func (c *DBClient) ShiftEvent(userID, eventID, activityID string, currentDate, newDate int64, event []byte) error {
	updates := []*badger.Entry{}
	deletes := [][]byte{}
	currentPrefix := []string{userKey, userID, eventKey, fmt.Sprint(currentDate), idKey, eventID, activityKey, activityID}
	deletes = append(deletes, key(currentPrefix))

	shiftedPrefix := []string{userKey, userID, eventKey, fmt.Sprint(newDate), idKey, eventID, activityKey, activityID}
	updates = append(updates, badger.NewEntry(key(shiftedPrefix), event))

	return updateDeleteItems(c, updates, deletes)
}

// Adds exercises instances to an event.
// The event must have been previously added.
// All existing exercise instances are deleted before the passed-in instances are added.
// The key of the exerciseInstances map is the index of the instances which determines order. The value is the instance.
// The key of the exerciseIDs map coincide with the index of the exerciseInstances map the index of the instances.
func (c *DBClient) AddExercisesToEvent(userID, eventID string, exerciseIDs map[int]string, exerciseInstances map[int][]byte) error {
	delPrefix := []string{userKey, userID, eventKey, eventID, exerciseKey}
	exEntries, err := readKeyPrefix(c, keyPrefix(delPrefix))
	if err != nil {
		return fmt.Errorf("error getting event exercises")
	}

	deletes := [][]byte{}
	for _, e := range exEntries {
		deletes = append(deletes, e.Key)
	}

	updates := []*badger.Entry{}

	addPrefix := []string{userKey, userID, eventKey, eventID, exerciseKey}

	for k, e := range exerciseInstances {
		exerciseID, ok := exerciseIDs[k]
		if !ok {
			return fmt.Errorf("mismatched key for exercise IDs")
		}
		prefix := append(addPrefix, exerciseID, indexKey, fmt.Sprint(k), instanceKey)
		entry := badger.NewEntry(key(prefix), e)
		updates = append(updates, entry)
	}

	if err := updateDeleteItems(c, updates, deletes); err != nil {
		return fmt.Errorf("failed to add exercie to event: %w", err)
	}

	return nil
}

// Returns a slice of exercise instances as byte slices.
func (c *DBClient) GetEventExercises(userID, eventID string) ([][]byte, error) {
	prefix := []string{userKey, userID, eventKey, eventID, exerciseKey}
	exEntries, err := readKeyPrefix(c, keyPrefix(prefix))
	if err != nil {
		return nil, fmt.Errorf("failed to read exercise entries: %w", err)
	}
	exInstancesrx := regexp.MustCompile(fmt.Sprintf("^%s:[^#]+#%s:[^#]+#%s:([^#]+)#%s:([^#]+)#%s$", userKey, eventKey, exerciseKey, indexKey, instanceKey))

	exercises := [][]byte{}

	for _, v := range exEntries {
		if exInstanceEntry := exInstancesrx.FindStringSubmatch(string(v.Key)); exInstanceEntry != nil {
			exercises = append(exercises, v.Value)

			continue
		}
	}

	return exercises, nil
}

// GetEvent retrieves an event from the database.
// Does not include the exercises
// If the event is not found nil is returned.
func (c *DBClient) GetEvent(userID, eventID string, eventDate int64) ([]byte, error) {
	prefix := []string{userKey, userID, eventKey, fmt.Sprint(eventDate), idKey, eventID, activityKey}
	entries, err := readKeyPrefix(c, keyPrefix(prefix))
	if err != nil {
		return nil, fmt.Errorf("failed to read event: %w", err)
	}
	if len(entries) == 0 {
		return nil, nil
	}
	event := entries[0].Value

	return event, nil
}

// GetEventPage gets a page of events.
// Items are iterated in reverse order so a value of zero for previous date returns no results.
// To get the latest events, set previous date to now.
// For subsequent pages, previousEventID and previousDate are used to identify the last item in the previous page
func (c *DBClient) GetEventPage(userID, previousEventID string, previousDate int64, pageSize int) (
	[][]byte, error) {
	// user:{id}#event:{date}#id:{id}#activity:{activityID}

	prefix := []string{userKey, userID, eventKey}
	startPrefix := make([]string, 3)
	firstPage := true
	var startKey []byte = nil
	if previousDate != 0 {
		copy(startPrefix, prefix)
		startPrefix = append(startPrefix, fmt.Sprint(previousDate))

		if previousEventID != "" {
			firstPage = false
			startPrefix = append(startPrefix, idKey, previousEventID)
		}
		startKey = key(startPrefix)
	}

	events := [][]byte{}

	entries, err := readKeyPrefixPage(c, startKey, key(prefix), pageSize, exerciseKey, firstPage, true)
	if err != nil {
		return nil, fmt.Errorf("failed to read events: %w", err)
	}

	for _, v := range entries {
		events = append(events, v.Value)
	}

	return events, nil
}

func (c *DBClient) GetKeys(usage string) (active, retired map[string][]byte, err error) {
	prefix := fmt.Sprintf("%s:", keyByCryptoUsage(usage))
	keyEntries, err := readKeyPrefix(c, []byte(prefix))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read keys: %w", err)
	}
	active = map[string][]byte{}
	retired = map[string][]byte{}

	rxp := regexp.MustCompile(fmt.Sprintf("^%s:([^#]+)$", keyByCryptoUsage(usage)))

	for _, v := range keyEntries {
		matches := rxp.FindStringSubmatch(string(v.Key))
		keyID := matches[1]
		if v.UserMeta == byte(1) {

			retired[keyID] = v.Value
		} else {
			active[keyID] = v.Value
		}
	}
	return active, retired, nil
}

// Adds a key of a certain type.
// Tags existing keys of that type with a meta flag.
func (c *DBClient) RotateKeys(newKey []byte, keyID, usage string) error {
	active, _, err := c.GetKeys(usage)
	if err != nil {
		return fmt.Errorf("failed to get existing keys: %w", err)
	}
	err = c.db.Update(func(txn *badger.Txn) error {
		for k, v := range active {
			prefix := []string{keyByCryptoUsage(usage), k}
			entry := badger.NewEntry(key(prefix), v).WithMeta(byte(1))

			if err := txn.SetEntry(entry); err != nil {
				return fmt.Errorf("failed to mark existing keys: %w", err)
			}
		}
		prefix := []string{keyByCryptoUsage(usage), keyID}
		newKeyEntry := badger.NewEntry(key(prefix), newKey).WithMeta(byte(0))
		if err := txn.SetEntry(newKeyEntry); err != nil {
			return fmt.Errorf("failed to add new key: %w", err)
		}
		return nil // func param
	})
	if err != nil {
		return fmt.Errorf("failed to rotate keys: %w", err)
	}
	return nil
}

func (c *DBClient) DeleteKey(keyID, usage string) error {
	keyPrefix := []string{keyByCryptoUsage(usage), keyID}
	keys := [][]byte{key(keyPrefix)}
	err := deleteKeys(c, keys)
	if err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}
	return nil
}

// Stores a user session.
// A session is stored across 3 items.
// ttl is in minutes
func (c *DBClient) AddSession(username, sessionID string, ttl int) error {
	prefix := []string{sessionKey, sessionID, userKey, username, expiresKey}

	expBytes, err := int64ToBSlice(time.Now().Add(time.Minute * time.Duration(ttl)).Unix())
	if err != nil {
		return fmt.Errorf("failed to convert expiry to byte slice: %w", err)
	}
	expEntry := badger.NewEntry(key(prefix), expBytes)
	updates := []*badger.Entry{expEntry}

	err = writeUpdates(c, updates)
	if err != nil {
		return fmt.Errorf("failed to add session: %w", err)
	}

	return nil
}

// Gets a session by ID.
func (c *DBClient) GetSession(sessionID string) (*string, *int64, error) {
	var username string
	var expiry *int64
	prefix := keyPrefix([]string{sessionKey, sessionID})
	sessionEntries, err := readKeyPrefix(c, prefix)
	if err != nil {
		return nil, nil, fmt.Errorf("could not read session: %w", err)
	}
	if len(sessionEntries) == 0 {
		return nil, nil, nil
	}
	e := sessionEntries[0]
	exp, err := bSliceToInt64(e.Value)
	if err != nil {
		return nil, nil, err
	}
	expiry = exp

	rxp := regexp.MustCompile(fmt.Sprintf(`^%s:%s#%s:([a-zA-Z0-9-]+)#%s`, sessionKey, sessionID, userKey, expiresKey))
	if sessKeyParts := rxp.FindStringSubmatch(string(e.Key)); sessKeyParts != nil {
		username = sessKeyParts[1]
	}

	return &username, expiry, nil
}

func (c *DBClient) DeleteSession(sessionID string) error {
	username, _, err := c.GetSession(sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	sessPrefix := []string{sessionKey, sessionID, userKey, *username, expiresKey}
	sessKey := key(sessPrefix)

	keys := [][]byte{sessKey}
	if err := deleteKeys(c, keys); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

// Returns a map of session IDs (keys) and expiry times.
func (c *DBClient) GetSessionExpiries() (map[string]int64, error) {
	prefix := []byte(sessionKey)
	entries, err := readKeyPrefix(c, prefix)
	if err != nil {
		return nil, fmt.Errorf("could not read sessions: %w", err)
	}
	expiries := map[string]int64{}
	for _, v := range entries {
		if strings.Contains(string(v.Key), expiresKey) {
			id := strings.Split(strings.Split(string(v.Key), "#")[0], ":")[1]
			exp, err := bSliceToInt64(v.Value)
			if err != nil {
				fmt.Println("binary.Read failed:", err)
			}
			expiries[id] = *exp
		}
	}

	return expiries, nil
}

func bSliceToInt64(b []byte) (*int64, error) {
	var exp int64
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &exp)
	if err != nil {
		return nil, err
	}
	return &exp, nil
}

func int64ToBSlice(i int64) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, i)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes(), nil
}

// Returns nil, nil if key not found
func readItem(c *DBClient, key []byte) (*badger.Entry, error) {
	var entry *badger.Entry
	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		value, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		entry = badger.NewEntry(item.KeyCopy(nil), value).WithMeta(item.UserMeta())
		return nil
	})
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return entry, nil
}

// func (c *DBClient) Iter8er() {
// 	c.db.View(func(txn *badger.Txn) error {
// 		opts := badger.DefaultIteratorOptions
// 		opts.PrefetchSize = 10
// 		it := txn.NewIterator(opts)
// 		defer it.Close()
// 		for it.seek(); it.Valid(); it.Next() {
// 			item := it.Item()
// 			k := item.Key()
// 			err := item.Value(func(v []byte) error {
// 				fmt.Printf("key=%s, value=%s\n", k, v)
// 				return nil
// 			})
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	})
// }

func readKeyPrefix(c *DBClient, prefix []byte) ([]*badger.Entry, error) {
	entries := []*badger.Entry{}
	itOptions := badger.DefaultIteratorOptions
	err := c.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(itOptions)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			entry := badger.NewEntry(item.KeyCopy(nil), val).WithMeta(item.UserMeta())
			entries = append(entries, entry)
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return entries, nil
}

func readKeyPrefixPage(c *DBClient, previousPrefix, validPrefix []byte, pageSize int, exclude string, includePrevious bool, reverse bool) ([]*badger.Entry, error) {
	entries := []*badger.Entry{}
	startPrefix := previousPrefix
	if previousPrefix == nil {
		startPrefix = validPrefix
	}
	itOptions := badger.DefaultIteratorOptions
	itOptions.Reverse = reverse
	err := c.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(itOptions)
		defer it.Close()
		for it.Seek(startPrefix); it.ValidForPrefix(validPrefix); it.Next() {
			item := it.Item()
			// Skip filtered items and the previous page's item (unless)
			if (previousPrefix != nil && strings.Contains(string(item.Key()), string(previousPrefix)) && !includePrevious) || (exclude != "" && strings.Contains(string(item.Key()), exclude)) {
				continue
			}
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			entry := badger.NewEntry(item.KeyCopy(nil), val)
			entries = append(entries, entry)

			if len(entries) == int(pageSize) {
				break
			}
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return entries, nil
}

// Deletes items and sets items in a single transaction
func updateDeleteItems(c *DBClient, updates []*badger.Entry, deletes [][]byte) error {
	err := c.db.Update(func(txn *badger.Txn) error {
		for _, v := range deletes {
			if err := txn.Delete(v); err != nil {
				return fmt.Errorf("delete failed: %w", err)
			}
		}
		for _, v := range updates {
			if err := txn.SetEntry(v); err != nil {
				return fmt.Errorf("transaction commit failed: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func writeUpdates(c *DBClient, updates []*badger.Entry) error {
	err := c.db.Update(func(txn *badger.Txn) error {
		for _, v := range updates {
			if err := txn.SetEntry(v); err != nil {
				return fmt.Errorf("transaction commit failed: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func deleteKeys(c *DBClient, keys [][]byte) error {
	err := c.db.Update(func(txn *badger.Txn) error {
		for _, v := range keys {
			if err := txn.Delete(v); err != nil {
				return fmt.Errorf("delete failed: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// use to check if a value that needs to be unique is not unique
// returns a pointer to the value if found, otherwise nil
func checkValue(c *DBClient, keyPrefix []byte, key, value string) (*string, error) {
	entries, err := readKeyPrefix(c, keyPrefix)
	if err != nil {
		return nil, err
	}
	for _, v := range entries {
		itemKey := strings.Split(string(v.Key), "#")[strings.Count(string(v.Key), "#")]
		if itemKey == key && bytes.Equal(v.Value, []byte(value)) {
			return &value, nil
		}
	}
	return nil, nil
}

func keyByCryptoUsage(usage string) string {
	switch usage {
	case "token":
		return tokenCryptoKeyKey
	}
	return ""
}

func key(items []string) []byte {
	leaf := ""
	if len(items)%2 == 1 {
		leaf = items[len(items)-1]
	}
	kvpairs := []string{}
	for i := 0; i < len(items)-1; i += 2 {
		kvpairs = append(kvpairs, fmt.Sprintf("%s:%s", items[i], items[i+1]))
	}
	key := strings.Join(kvpairs, "#")
	if leaf != "" {
		key = fmt.Sprintf("%s#%s", key, leaf)
	}
	return []byte(key)
}

func keyPrefix(items []string) []byte {
	key := string(key(items))
	if len(items)%2 == 0 {
		return []byte(key + "#")
	}
	return []byte(key + ":")
}
