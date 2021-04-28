package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

type Evoting struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
type Asset struct {
	ID             string `json:"ID"`
	Color          string `json:"color"`
	Size           int    `json:"size"`
	Owner          string `json:"owner"`
	AppraisedValue int    `json:"appraisedValue"`
}

type User struct {
	ID       string
	Name     string
	Email    string
	Password string
	Doctype  string
}

type Quote struct {
	Value string
}

type Voter struct {
	ID         string
	Name       string
	Email      string
	Vote       string
	ElectionID string
	Doctype    string
}

type Election struct {
	ID      string
	Name    string
	Doctype string
	Ended   bool
}

type Candidate struct {
	ID         string
	Name       string
	ElectionID string
	Marka      string
	Doctype    string
}

type Vote struct {
	ID          string
	ElectionID  string
	CandidateID string
	Doctype     string
}

type ElectionResult struct {
	ElectionID  string
	CandidateID string
	Marka       string
	VoteCount   int64
}

func (s *SmartContract) CalculateResult(ctx contractapi.TransactionContextInterface, electionid string) ([]*ElectionResult, error) {

	queryString := newCouchQueryBuilder().addSelector("Doctype", "vote").addSelector("ElectionID", electionid).getQueryString()

	fmt.Println(queryString)

	iterator, _ := ctx.GetStub().GetQueryResult(queryString)
	defer iterator.Close()

	var votes []*Vote
	for iterator.HasNext() {
		queryResult, err := iterator.Next()
		if err != nil {
			return nil, err
		}
		var vote Vote
		err = json.Unmarshal(queryResult.Value, &vote)
		if err != nil {
			return nil, err
		}
		votes = append(votes, &vote)
	}

	voteCount := make(map[string]int64)

	for _, vote := range votes {
		voteCount[vote.CandidateID] += 1
	}

	result := make([]*ElectionResult, 0)

	for candidateId, numOfVotes := range voteCount {
		candidateJSON, _ := ctx.GetStub().GetState(candidateId)
		var candidate Candidate
		json.Unmarshal(candidateJSON, &candidate)

		candidateResult := ElectionResult{
			ElectionID:  electionid,
			CandidateID: candidateId,
			Marka:       candidate.Marka,
			VoteCount:   numOfVotes,
		}
		result = append(result, &candidateResult)
	}

	return result, nil

}

// func (s *SmartContract) totalCandidates(electionid string, stub shim.ChaincodeStubInterface) int64 {
// 	queryString := newCouchQueryBuilder().addSelector("doctype", "Candidate").addSelector("electionid", electionid).getQueryString()

// 	iterator, _ := stub.GetQueryResult(queryString)
// 	counter := 0

// 	for iterator.HasNext() {
// 		counter++
// 		resp, _ := iterator.Next()
// 		fmt.Println(string(resp.Value))
// 	}
// 	return int64(counter)
// }

////////have to do testing//////////////////
func (s *SmartContract) VoteCasting(ctx contractapi.TransactionContextInterface, id string, electionid string, candidateid string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	asset := Vote{
		ID:          id,
		ElectionID:  electionid,
		CandidateID: candidateid,
		Doctype:     "vote",
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

func (s *SmartContract) AddCandidate(ctx contractapi.TransactionContextInterface, id string, name string, marka string, electionid string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	asset := Candidate{
		ID:         id,
		Name:       name,
		Marka:      marka,
		ElectionID: electionid,
		Doctype:    "candidate",
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

func getUsersQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*User, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructUsersQueryResponseFromIterator(resultsIterator)
}

//////////////////

func (t *SmartContract) ShowAllElections(ctx contractapi.TransactionContextInterface) ([]*Election, error) {
	queryString := newCouchQueryBuilder().addSelector("Doctype", "election").getQueryString() // fmt.Sprintf(`{"selector":{"DocType":"user","Email","Password":"%s","%s"}}`, email, password)
	fmt.Println(queryString)

	elections, err := getElectionsQueryResultForQueryString(ctx, queryString)

	if err != nil {
		return nil, err
	}

	return elections, nil
}
func getElectionsQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Election, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructElectionsQueryResponseFromIterator(resultsIterator)
}

func constructElectionsQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*Election, error) {
	var elections []*Election
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var election Election
		err = json.Unmarshal(queryResult.Value, &election)
		if err != nil {
			return nil, err
		}
		elections = append(elections, &election)
	}

	return elections, nil
}

///////////
func (t *SmartContract) ShowAllCandidates(ctx contractapi.TransactionContextInterface, electionid string) ([]*Candidate, error) {
	queryString := newCouchQueryBuilder().addSelector("Doctype", "candidate").addSelector("ElectionID", electionid).getQueryString() // fmt.Sprintf(`{"selector":{"DocType":"user","Email","Password":"%s","%s"}}`, email, password)
	fmt.Println(queryString)

	candidates, err := getCandidatesQueryResultForQueryString(ctx, queryString)

	if err != nil {
		return nil, err
	}

	return candidates, nil
}

func getCandidatesQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Candidate, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructCandidatesQueryResponseFromIterator(resultsIterator)
}
func constructCandidatesQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*Candidate, error) {
	var candidates []*Candidate
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var candidate Candidate
		err = json.Unmarshal(queryResult.Value, &candidate)
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, &candidate)
	}

	return candidates, nil
}

func constructUsersQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*User, error) {
	var users []*User
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var user User
		err = json.Unmarshal(queryResult.Value, &user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{ID: "asset1", Color: "blue", Size: 5, Owner: "Tomoko", AppraisedValue: 300},
		{ID: "asset2", Color: "red", Size: 5, Owner: "Brad", AppraisedValue: 400},
		{ID: "asset3", Color: "green", Size: 10, Owner: "Jin Soo", AppraisedValue: 500},
		{ID: "asset4", Color: "yellow", Size: 10, Owner: "Max", AppraisedValue: 600},
		{ID: "asset5", Color: "black", Size: 15, Owner: "Adriana", AppraisedValue: 700},
		{ID: "asset6", Color: "white", Size: 15, Owner: "Michel", AppraisedValue: 800},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	asset := Asset{
		ID:             id,
		Color:          color,
		Size:           size,
		Owner:          owner,
		AppraisedValue: appraisedValue,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

/////////////////todo
func (s *SmartContract) CreateElection(ctx contractapi.TransactionContextInterface, id string, name string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	asset := Election{
		ID:      id,
		Name:    name,
		Ended:   false,
		Doctype: "election",
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

func (s *SmartContract) StopElection(ctx contractapi.TransactionContextInterface, id string) error {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("failed to stop election: %v", err)
	}

	var election Election
	err = json.Unmarshal(assetJSON, &election)

	if err != nil {
		return fmt.Errorf("failed to stop election: %v", err)
	}

	election.Ended = true

	assetJSON, err = json.Marshal(election)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

func (t *SmartContract) LoginUser(ctx contractapi.TransactionContextInterface, email string, password string) ([]*User, error) {
	queryString := newCouchQueryBuilder().addSelector("Doctype", "user").addSelector("Email", email).addSelector("Password", password).getQueryString() // fmt.Sprintf(`{"selector":{"DocType":"user","Email","Password":"%s","%s"}}`, email, password)
	fmt.Println(queryString)

	users, err := getUsersQueryResultForQueryString(ctx, queryString)

	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("email or password incorrect")
	}

	return users, nil
}

//////////////working////////////
func (s *SmartContract) RegisterUser(ctx contractapi.TransactionContextInterface, id string, name string, email string, password string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	asset := User{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: password,
		Doctype:  "user",
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

func (s *SmartContract) SayHello(ctx contractapi.TransactionContextInterface) (*Quote, error) {
	quote := Quote{"Hello World!"}

	return &quote, nil
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	// overwriting original asset with new asset
	asset := Asset{
		ID:             id,
		Color:          color,
		Size:           size,
		Owner:          owner,
		AppraisedValue: appraisedValue,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return err
	}

	asset.Owner = newOwner
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}
