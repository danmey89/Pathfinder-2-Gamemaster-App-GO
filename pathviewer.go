package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var base = map[string]string{
	"Perception": "Wis", "Fortitude": "Con", "Reflex": "Dex", "Will": "Wis", "Acrobatics": "Dex", "Arcana": "Int", "Athletics": "Str",
	"Crafting": "Int", "Deception": "Cha", "Diplomacy": "Cha", "Intimidation": "Cha", "Medicine": "Wis", "Nature": "Wis", "Occultism": "Int",
	"Performance": "Cha", "Religion": "Wis", "Society": "Int", "Stealth": "Dex", "Survival": "Wis", "Thievery": "Dex"}

var abilities = [6]string{"Str", "Dex", "Con", "Int", "Wis", "Cha"}

var proficiencies = [20]string{"perception", "fortitude", "reflex", "will", "acrobatics", "arcana", "athletics", "crafting", "deception", "diplomacy",
	"intimidation", "medicine", "nature", "occultism", "performance", "religion", "society", "stealth", "survival", "thievery"}

var flagVar1 bool
var flagVar2 bool
var charMap = make(map[int]map[string]interface{})

var templates = template.Must(template.ParseFiles("templates/index.html"))

var fs = http.FileServer(http.Dir("./static"))

func init() {
	flag.BoolVar(&flagVar1, "createDB", false, "Initialize Database")
	flag.BoolVar(&flagVar2, "saveData", false, "Load JSON into Database")
}

func main() {

	flag.Parse()

	if flagVar1 {
		createDatabase()
	}

	pathfinderDB, err := sql.Open("sqlite3", "pathfinder.sqlite")
	if err != nil {
		log.Println(err)
	}
	defer pathfinderDB.Close()

	if flagVar2 {
		saveData(pathfinderDB)
	}
	if !flagVar1 && !flagVar2 {
		loadCharacters(pathfinderDB)

		mux := http.NewServeMux()

		mux.HandleFunc("/", indexHandler)
		mux.Handle("/static/", http.StripPrefix("/static", fs))

		fmt.Println("Server running at http://localhost:8080")
		log.Fatal(http.ListenAndServe(":8080", mux))
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	p := Page{
		Proficiencies: proficiencies,
		Abilities:     abilities,
		Data:          charMap,
	}

	err := templates.ExecuteTemplate(w, "index.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func saveData(db *sql.DB) {

	dir, err := os.ReadDir("./docs")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range dir {

		path := filepath.Join("docs", file.Name())

		c := getCharacter(path)

		insertRow(&c, db)
	}
}

func getCharacter(path string) Character {

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Succesfully opened File")
	}

	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var char CharacterWrap

	json.Unmarshal(byteValue, &char)

	setTraining(&char.Character)
	setModifiers(&char.Character)
	calculateModifiers(&char.Character)

	log.Printf("%s has been imported \n", char.Character.Name)

	return char.Character

}

func setTraining(p *Character) {

	pro := p.Proficiencies
	train := &p.Training
	s1 := reflect.ValueOf(pro)

	for k := range base {
		p := reflect.Indirect(s1).FieldByName(k).Int()

		var t string

		switch p {
		case 0:
			t = "u"
		case 2:
			t = "t"
		case 4:
			t = "e"
		case 6:
			t = "m"
		case 8:
			t = "l"
		}
		reflect.ValueOf(train).Elem().FieldByName(k).SetString(t)
	}
}

func calculateModifiers(p *Character) {

	pro := &p.Proficiencies

	s0 := reflect.ValueOf(p.AbilityModifier)
	s1 := reflect.ValueOf(pro)
	s2 := reflect.ValueOf(p.Level)
	l := reflect.Indirect(s2).Int()

	for k, v := range base {

		a := reflect.Indirect(s0).FieldByName(v).Int()
		p := reflect.Indirect(s1).FieldByName(k).Int()

		sum := l + a + p

		reflect.ValueOf(pro).Elem().FieldByName(k).SetInt(sum)
	}

}

func setModifiers(c *Character) {

	for _, m := range abilities {

		s := reflect.ValueOf(c.Abilities)
		f := reflect.Indirect(s).FieldByName(m).Int()

		r := (f - 10) / 2

		reflect.ValueOf(&c.AbilityModifier).Elem().FieldByName(m).SetInt(r)
	}
}

func createDatabase() {

	os.Remove("pathfinder.sqlite")

	log.Println("Creating Database")

	_, err := os.Create("pathfinder.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database Created")

	sqlCmd, err := os.ReadFile("schema.sql")
	if err != nil {
		log.Fatal(err)
	}
	schema := string(sqlCmd)

	pathfinderDB, err := sql.Open("sqlite3", "./pathfinder.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer pathfinderDB.Close()

	pathfinderDB.Exec(schema)

}

func insertRow(c *Character, db *sql.DB) {

	name := c.Name

	log.Printf("Insert Character %s\n", name)

	in := []string{"Name", "Class", "Level", "Ancestry", "Heritage", "Background", "AC", "Str", "Dex", "Con", "Int", "Wis", "Cha",
		"StrMod", "DexMod", "ConMod", "IntMod", "WisMod", "ChaMod",
		"Perception", "Fortitude", "Reflex", "Will", "Acrobatics", "Arcana", "Athletics", "Crafting", "Deception", "Diplomacy",
		"Intimidation", "Medicine", "Nature", "Occultism", "Performance", "Religion", "Society", "Stealth", "Survival", "Thievery",
		"PerceptionTrain", "FortitudeTrain", "ReflexTrain", "WillTrain", "AcrobaticsTrain", "ArcanaTrain", "AthleticsTrain", "CraftingTrain",
		"DeceptionTrain", "DiplomacyTrain", "IntimidationTrain", "MedicineTrain", "NatureTrain", "OccultismTrain", "PerformanceTrain",
		"ReligionTrain", "SocietyTrain", "StealthTrain", "SurvivalTrain", "ThieveryTrain", "Languages"}

	var q = make([]string, len(in))

	for i := range q {
		q[i] = "?"
	}

	insertChar := fmt.Sprintf(`INSERT INTO characters (%s) VALUES (%s)`, strings.Join(in, ", "), strings.Join(q, ", "))

	statement, err := db.Prepare(insertChar)
	if err != nil {
		log.Fatal(err)
	}

	_, err = statement.Exec(c.Name, c.Class, c.Level, c.Ancestry, c.Heritage, c.Background, c.AC.AC, c.Abilities.Str, c.Abilities.Dex, c.Abilities.Con,
		c.Abilities.Int, c.Abilities.Wis, c.Abilities.Cha, c.AbilityModifier.Str, c.AbilityModifier.Dex, c.AbilityModifier.Con, c.AbilityModifier.Int,
		c.AbilityModifier.Wis, c.AbilityModifier.Cha, c.Proficiencies.Perception, c.Proficiencies.Fortitude, c.Proficiencies.Reflex, c.Proficiencies.Will,
		c.Proficiencies.Acrobatics, c.Proficiencies.Arcana, c.Proficiencies.Athletics, c.Proficiencies.Crafting, c.Proficiencies.Deception, c.Proficiencies.Diplomacy,
		c.Proficiencies.Intimidation, c.Proficiencies.Medicine, c.Proficiencies.Nature, c.Proficiencies.Occultism, c.Proficiencies.Performance, c.Proficiencies.Religion,
		c.Proficiencies.Society, c.Proficiencies.Stealth, c.Proficiencies.Survival, c.Proficiencies.Thievery, c.Training.Perception, c.Training.Fortitude, c.Training.Reflex,
		c.Training.Will, c.Training.Acrobatics, c.Training.Arcana, c.Training.Athletics, c.Training.Crafting, c.Training.Deception, c.Training.Diplomacy,
		c.Training.Intimidation, c.Training.Medicine, c.Training.Nature, c.Training.Occultism, c.Training.Performance, c.Training.Religion, c.Training.Society,
		c.Training.Stealth, c.Training.Survival, c.Training.Thievery, strings.Join(c.Languages, ", "))
	if err != nil {
		log.Println(err)
	} else {
		log.Println("\n has been inserted")
	}

}

func loadCharacters(db *sql.DB) map[int]map[string]interface{} {

	rows, err := db.Query(`SELECT * FROM characters`)
	if err != nil {
		log.Fatal(err)
	}
	cols, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	output := make(map[int]map[string]interface{})
	i := 0

	for rows.Next() {

		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			log.Fatal(err)
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}

		charMap[i] = m
		i++
	}
	return output
}

type CharacterWrap struct {
	Exists    bool      `json:"success"`
	Character Character `json:"build"`
}

type Character struct {
	Name            string          `json:"name"`
	Class           string          `json:"class"`
	Level           int             `json:"level"`
	Ancestry        string          `json:"ancestry"`
	Heritage        string          `json:"heritage"`
	Background      string          `json:"background"`
	Languages       []string        `json:"languages"`
	Abilities       Abilities       `json:"abilities"`
	Proficiencies   Proficiencies   `json:"proficiencies"`
	Lore            [][]interface{} `json:"lores"`
	AC              AC              `json:"acTotal"`
	Training        Training
	AbilityModifier Abilities
}

type Abilities struct {
	Str int `json:"str"`
	Dex int `json:"dex"`
	Con int `json:"con"`
	Int int `json:"int"`
	Wis int `json:"wis"`
	Cha int `json:"cha"`
}

type Proficiencies struct {
	Perception   int `json:"perception"`
	Fortitude    int `json:"fortitude"`
	Reflex       int `json:"reflex"`
	Will         int `json:"will"`
	Acrobatics   int `json:"acrobatics"`
	Arcana       int `json:"arcana"`
	Athletics    int `json:"athletics"`
	Crafting     int `json:"crafting"`
	Deception    int `json:"deception"`
	Diplomacy    int `json:"diplomacy"`
	Intimidation int `json:"intimidation"`
	Medicine     int `json:"medicine"`
	Nature       int `json:"nature"`
	Occultism    int `json:"occultism"`
	Performance  int `json:"performance"`
	Religion     int `json:"religion"`
	Society      int `json:"society"`
	Stealth      int `json:"stealth"`
	Survival     int `json:"survival"`
	Thievery     int `json:"thievery"`
}

type Training struct {
	Perception   string
	Fortitude    string
	Reflex       string
	Will         string
	Acrobatics   string
	Arcana       string
	Athletics    string
	Crafting     string
	Deception    string
	Diplomacy    string
	Intimidation string
	Medicine     string
	Nature       string
	Occultism    string
	Performance  string
	Religion     string
	Society      string
	Stealth      string
	Survival     string
	Thievery     string
}

type AC struct {
	AC int `json:"acTotal"`
}

type Page struct {
	Proficiencies [20]string
	Abilities     [6]string
	Data          map[int]map[string]interface{}
}
