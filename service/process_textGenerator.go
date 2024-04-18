package service

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/milla-v/xandr/bss"
	"github.com/milla-v/xandr/bss/xgen"
)

func checkSeparators(seps separators) error {
	notAllowed := "[](){}$\\/|?*+-"
	sp := []string{
		seps.Sep1, seps.Sep2, seps.Sep3, seps.Sep4, seps.Sep5,
	}

	for i, s := range sp {
		log.Println(i, ". ", s)
		log.Println("len(s) = ", len(s))
		if len(s) != 1 && s != "TAB" && s != "SPACE" {
			return fmt.Errorf("sep%d should be a single character", i+1)
		}
		if s != "TAB" && s != "SPACE" {
			log.Println("s != tab or space: ", s)
		}
		if strings.ContainsAny(s, notAllowed) {
			return fmt.Errorf("sep%d: symbols "+notAllowed+" are not allowed as a separators", i+1)
		}
	}
	return nil
}

func setDefaultSeparators(seps *separators) {
	if seps.Sep1 == "" {
		seps.Sep1 = ":"
	}
	if seps.Sep2 == "" {
		seps.Sep2 = ";"
	}
	if seps.Sep3 == "" {
		seps.Sep3 = ","
	}
	if seps.Sep4 == "" {
		seps.Sep4 = "#"
	}
	if seps.Sep5 == "" {
		seps.Sep5 = "^"
	}
}

func generateSegments(segmentFields []string, seps separators, count int) (segmentsToAdd, segmentsToRem []string) {
	for i := 0; i < count; i++ {
		var segmentAdd []string
		var segmentRem []string

		for _, sf := range segmentFields {
			switch sf {
			case "SEG_ID":
				segID := 1000 + rand.Intn(1000)
				segmentAdd = append(segmentAdd, strconv.Itoa(segID))
				segmentRem = append(segmentRem, strconv.Itoa(segID+100))
			case "SEG_CODE":
				segID := 1000 + rand.Intn(1000)
				segmentAdd = append(segmentAdd, "code_"+strconv.Itoa(segID))
				segmentRem = append(segmentRem, "code_"+strconv.Itoa(segID+100))
			case "MEMBER_ID":
				segmentAdd = append(segmentAdd, "100")
				segmentRem = append(segmentRem, "100")
			case "EXPIRATION":
				segmentAdd = append(segmentAdd, "43200")
				segmentRem = append(segmentRem, "-1")
			case "VALUE":
				value := 1 + rand.Intn(5)
				segmentAdd = append(segmentAdd, strconv.Itoa(value))
				segmentRem = append(segmentRem, "0")
			case "TIMESTAMP":
				segID := 1000 + rand.Intn(1000)
				segmentAdd = append(segmentAdd, strconv.Itoa(segID))
				segmentRem = append(segmentRem, strconv.Itoa(segID+100))
			}

		}

		segmentsToAdd = append(segmentsToAdd, strings.Join(segmentAdd, seps.Sep3))
		segmentsToRem = append(segmentsToRem, strings.Join(segmentRem, seps.Sep3))
	}

	return
}

func replaceTabs(s string) string {
	if s == "TAB" {
		return "\t"
	}
	if s == "SPACE" {
		return " "
	}
	return s
}

/*
func generateSample2(segmentFields []string, seps separators) string {
	const lineTemplate = "{UID}{SEP_1}{SEGMENTS_ADD}{SEP_4}{SEGMENTS_DEL}{SEP_5}{DOMAIN}"

	idtypes := []idtype{
		{"xandrid", 0},
		{"idfa", 3},
		{"aaid", 8},
	}

	seps.Sep1 = replaceTabs(seps.Sep1)
	seps.Sep4 = replaceTabs(seps.Sep4)

	var s string

	for _, idt := range idtypes {

		var uid string

		switch idt.domain {
		case "xandrid":
			uid = strconv.Itoa(int(rand.Int63()))
		case "idfa", "aaid":
			uid = guuid.New().String()
		default:
			log.Println("ERROR: invalid domain", idt.domain)
			continue
		}

		segmentsToAdd, segmentsToRem := generateSegments(segmentFields, seps, 2)

		var domain, sep5 string
		if idt.number != 0 {
			domain = strconv.Itoa(idt.number)
			sep5 = seps.Sep5
		}

		sr := strings.NewReplacer(
			"{UID}", uid,
			"{SEP_1}", seps.Sep1,
			"{SEGMENTS_ADD}", strings.Join(segmentsToAdd, seps.Sep2),
			"{SEP_4}", seps.Sep4,
			"{SEGMENTS_DEL}", strings.Join(segmentsToRem, seps.Sep2),
			"{SEP_5}", sep5,
			"{DOMAIN}", domain,
		)
		log.Println("sr:", sr)
		s += sr.Replace(lineTemplate) + "\n"
		log.Println("sr:", sr)
	}

	return s
}

*/

func generateSample(params *xgen.TextEncoderParameters) (string, error) {
	var err error
	var out bytes.Buffer

	var users []*xgen.UserRecord

	user1 := &xgen.UserRecord{
		UID: "12345",
		Segments: []xgen.Segment{
			{ID: 100, Expiration: 180 * 60 * 24, Value: 123, MemberID: 700, Code: "code_123", Timestamp: time.Now().Unix()},
			{ID: 101, Expiration: 1440, Value: 123, MemberID: 700, Code: "code_123", Timestamp: time.Now().Unix()},
		},
	}

	user2 := &xgen.UserRecord{
		UID:    "12345",
		Domain: xgen.IDFA,
		Segments: []xgen.Segment{
			{ID: 100, Expiration: 180 * 60 * 24, Value: 123, MemberID: 700, Code: "code_123", Timestamp: time.Now().Unix()},
			{ID: 101, Expiration: 1440, Value: 123, MemberID: 700, Code: "code_123", Timestamp: time.Now().Unix()},
		},
	}

	log.Println("User1: ", user1)
	log.Println("Timestemp for user1: ", user1.Segments[0].Timestamp, " | ", user1.Segments[1].Timestamp)
	log.Println("------------------")
	log.Println("User2: ", user2)
	log.Println("Timestemp for user2: ", user2.Segments[0].Timestamp, " | ", user2.Segments[1].Timestamp)

	users = append(users, user1, user2)

	w, err := bss.NewSegmentDataFormatter(&out, bss.FormatText, params)
	if err != nil {
		log.Println(err)
		return "", err
	}

	if err := w.Append(users); err != nil {
		log.Println(err)
		return "", err
	}

	if err := w.Close(); err != nil {
		log.Println(err)
		return "", err
	}

	return out.String(), err
}

func checkSegments(segmentFields []string) (string, error) {
	var err error
	var check string
	var segIDfound bool
	var segCodeFound bool

	//start check segmentFields
	for _, s := range segmentFields {
		if strings.Contains(s, "SEG_ID") {
			segIDfound = true
		}
		if strings.Contains(s, "SEG_CODE") {
			segCodeFound = true
		}
	}
	//check if at least  SEG_ID or SEG_CODE was choosen
	if segIDfound == false && segCodeFound == false {
		check = "Choose at least  SEG_ID or SEG_CODE"
	}
	// check if SEG_CODE or SEG_ID included but not both.
	if segIDfound == true && segCodeFound == true {
		check = "You may include SEG_CODE or SEG_ID but not both."
	}

	return check, err
}
