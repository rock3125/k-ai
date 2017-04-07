package freebase

import (
	"k-ai/db"
	"errors"
	"bytes"
	"k-ai/util"
	"strings"
	"strconv"
	"k-ai/logger"
	"k-ai/nlu/model"
)

type FBIndex struct {
	Predicate   int     `json:"predicate"`
	Word        int     `json:"word"`
	Tuples      []int   `json:"tuples"`
}

type FBTuple struct {
	Id        int       `json:"id"`
	Lhs       []int     `json:"lhs"`
	Predicate   int     `json:"predicate"`
	Rhs       []int     `json:"rhs"`
}

// the string version of FBTuple
type FBTupleString struct {
	Id        int
	Lhs       []string
	Predicate   string
	Rhs       []string
}


func GetExamples() error {
	content, err := util.LoadTextFile("/home/peter/dev/list.txt")
	if err != nil {
		return err
	}

	int_list := make([]int,0)
	parts := strings.Split(content, ",")
	for _, part := range parts {
		value, err := strconv.Atoi(part)
		if err != nil {
			return err
		}
		int_list = append(int_list,value)
	}

	tuple_list, err := getTuplesFromIdList(int_list)
	if err != nil {
		return err
	}

	tuple_string_list, err := tupleListToTupleStringList(tuple_list)
	if err != nil {
		return err
	}

	for _, tuple := range tuple_string_list {
		logger.Log.Info(tuple.String())
	}

	return nil
}


// convert the contents of an FBTupleString to string
func (tuple FBTupleString) String() string {
	var buffer bytes.Buffer
	if tuple.Lhs != nil {
		for _, str := range tuple.Lhs {
			buffer.WriteString(str)
			buffer.WriteString(" ")
		}
	}
	buffer.WriteString(tuple.Predicate)
	buffer.WriteString(" ")
	if tuple.Rhs != nil {
		for _, str := range tuple.Rhs {
			buffer.WriteString(str)
			buffer.WriteString(" ")
		}
	}
	return buffer.String()
}

// get an index by predicate and word id
func (index *FBIndex) Get() error {
	if index.Word == 0 || index.Predicate == 0 {
		return errors.New("invalid parameters")
	}
	columns := []string{"tuples"}
	whereMap := make(map[string]interface{},0)
	whereMap["predicate"] = index.Predicate
	whereMap["word"] = index.Word

	select_str := db.Cassandra.SelectPaginated("freebase_index", columns,
		whereMap, "", nil, 0)

	iter := db.Cassandra.Session.Query(select_str).Iter()
	var int_list []int
	if iter.Scan(&int_list) {
		index.Tuples = int_list
	}
	return iter.Close()
}

// get a single tuple by ID
func (tuple *FBTuple) Get() error {
	if tuple.Id == 0  {
		return errors.New("invalid parameters")
	}
	columns := []string{"lhs", "predicate", "rhs"}
	whereMap := make(map[string]interface{},0)
	whereMap["id"] = tuple.Id

	select_str := db.Cassandra.SelectPaginated("freebase_tuple", columns,
		whereMap, "", nil, 0)

	iter := db.Cassandra.Session.Query(select_str).Iter()
	var predicate int
	var lhs []int
	var rhs []int
	if iter.Scan(&lhs, &predicate, &rhs) {
		tuple.Predicate = predicate
		tuple.Lhs = lhs
		tuple.Rhs = rhs
	}
	return iter.Close()
}

// convert an int tuple to a string based tuple
func (tuple FBTuple) ToString() (*FBTupleString, error) {
	string_int_list := make([]int,0)
	for _, id := range tuple.Lhs {
		string_int_list = append(string_int_list, id)
	}
	string_int_list = append(string_int_list, tuple.Predicate)
	for _, id := range tuple.Rhs {
		string_int_list = append(string_int_list, id)
	}
	word_id_list, err := db.Cassandra.FreebaseIdsToStringList(string_int_list)
	if err != nil {
		return nil, err
	}
	// setup a fast lookup map
	word_id_map := make(map[int]string,0)
	for _, word_id := range word_id_list {
		word_id_map[word_id.Id] = word_id.Word
	}
	result := FBTupleString{Id: tuple.Id, Lhs: make([]string,0), Rhs: make([]string,0)}
	for _, id := range tuple.Lhs {
		result.Lhs = append(result.Lhs, word_id_map[id])
	}
	for _, id := range tuple.Rhs {
		result.Rhs = append(result.Rhs, word_id_map[id])
	}
	result.Predicate = word_id_map[tuple.Predicate]
	return &result, nil
}

// get a list of tuples from the ids provided
func getTuplesFromIdList(id_list []int) ([]FBTuple, error) {
	if len(id_list) == 0 {
		return nil, errors.New("invalid parameters")
	}
	tuple_list := make([]FBTuple,0)
	for _, id := range id_list {
		tuple := FBTuple{Id: id}
		err := tuple.Get()
		if err != nil {
			return nil, err
		}
		tuple_list = append(tuple_list, tuple)
	}
	return tuple_list, nil
}


// get a list of tuples from the ids provided
func tupleListToTupleStringList(tuple_list []FBTuple) ([]*FBTupleString, error) {
	if len(tuple_list) == 0 {
		return nil, errors.New("invalid parameters")
	}
	tuple_string_list := make([]*FBTupleString,0)
	for _, tuple := range tuple_list {
		tuple_string, err := tuple.ToString()
		if err != nil {
			return nil, err
		}
		tuple_string_list = append(tuple_string_list, tuple_string)
	}
	return tuple_string_list, nil
}


// perform a search for the search terms specified across freebase
// returns a list of Clause Ids and a list of the words that matched, or error
func freebaseFind(terms []string, page int, page_size int) ([]int, []db.WordId, error) {
	word_id_list, err := db.Cassandra.FreebaseStringsToIdList(terms)
	if err != nil {
		return nil, nil, err
	}
	// we must have at least one predicate, and one term to search
	if len(word_id_list) <= 1 {
		return nil, nil, errors.New("term list not complete or empty")
	}
	// must have one predicate
	num_predicates := 0
	predicate := db.WordId{}
	for _, word_id := range word_id_list {
		if word_id.Is_predicate {
			predicate.Id = word_id.Id
			predicate.Word = word_id.Word
			predicate.Is_predicate = true
			num_predicates += 1
		}
	}
	if num_predicates == 0 {
		return nil, nil, errors.New("no predicates in term list")
	}
	if num_predicates > 1 {
		return nil, nil, errors.New("too many predicates in term list")
	}

	// perform the search - get the indexes and intersect them to find a suitable tuple
	intersection_set := make(map[int]int, 0)
	hit_words := make([]db.WordId,0)
	hit_words = append(hit_words, predicate)
	counter := 0
	for _, word_id := range word_id_list {
		if word_id.Id != predicate.Id {
			index_search := FBIndex{Predicate: predicate.Id, Word: word_id.Id}
			err := index_search.Get()
			if err == nil { // no error, go ahead
				// any words associated with this one?
				if len(index_search.Tuples) > 0 {
					if counter == 0 { // first time around, take all items
						hit_words = append(hit_words, word_id)
						for _, id := range index_search.Tuples {
							intersection_set[id] = id
						}
					} else {
						// all other times, create a new map each time with surviving members
						new_intersection_set := make(map[int]int, 0)
						for _, id := range index_search.Tuples {
							if val, ok := intersection_set[id]; ok {
								new_intersection_set[val] = val
							}
						}
						// move map to old map, and exit if empty
						intersection_set = new_intersection_set
						if len(new_intersection_set) == 0 {
							break
						}
						hit_words = append(hit_words, word_id) // collect words that contributed
					}
					counter += 1 // next iteration / keyword
				} // if index has words associated with it
			}
		} // if word not "the" predicate
	} // for each word

	// empty?  return empty lists and no error
	result_list := make([]int, 0)
	if len(intersection_set) == 0 {
		hit_words = make([]db.WordId,0) // empty it
		return result_list, hit_words, nil
	}
	// otherwise, create the result list
	counter = 0
	start_page := page * page_size
	end_page := start_page + page_size
	for key, _ := range intersection_set {
		if start_page <= counter && counter < end_page {
			result_list = append(result_list, key)
		}
		counter += 1
	}
	return result_list, hit_words, nil
}

// query Freebase
func FreebaseQuery(query_list []string, page int, page_size int) ([]*FBTupleString, error) {
	tuple_id_list, _, err := freebaseFind(query_list, page, page_size)
	if err == nil {
		tuple_list, err := getTuplesFromIdList(tuple_id_list)
		if err == nil {
			tuple_string_list, err := tupleListToTupleStringList(tuple_list)
			if err == nil {
				return tuple_string_list, err
			}
		}
	}
	return nil, err
}


// search Freebase (well - use the freebaseSearch structure to query)
func FreebaseQueryBySearch(token_list []model.Token, page int, page_size int) ([]*FBTupleString, error) {
	search, err := MatchSystem.Match(token_list)
	if err != nil {
		return nil, err
	}
	query_list := make([]string,0)
	query_list = append(query_list, search.Predicate)
	for _, t_token := range search.TokenList {
		query_list = append(query_list, t_token.Text)
	}
	tuple_id_list, _, err := freebaseFind(query_list, page, page_size)
	if err == nil {
		tuple_list, err := getTuplesFromIdList(tuple_id_list)
		if err == nil {
			tuple_string_list, err := tupleListToTupleStringList(tuple_list)
			if err == nil {
				return tuple_string_list, err
			}
		}
	}
	return nil, err
}

