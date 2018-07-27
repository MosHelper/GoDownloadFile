package uploadtodrive

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type driveResp struct {
	Result  int    `json:"result"`
	URL     string `json:"url"`
	Convert bool   `json:"convert"`
}

type driveFile struct {
	WebContentLink string `json:"webContentLink"`
}

func Upload(pathfile string, filename string) (string, error) {
	// read file for upload purpose
	fileToBeChunked := pathfile

	file, err := os.Open(fileToBeChunked)

	if err != nil {
		log.Println(err)
		return "", err
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()

	log.Println("Totle Size : ", strconv.FormatInt(fileSize, 10))

	/////////////////////////////////////////////////////////////////
	url := "https://sdrive.online/wp-admin/admin-ajax.php"

	payload := strings.NewReader(`action=useyourdrive-upload-file&type=get-direct-url&filename=` + filename + `&file_size=` + strconv.FormatInt(fileSize, 10) + `&mimetype=application%2Fzip&orgin=https%3A%2F%2Fsdrive.online&lastFolder=1OnEI9_QOoGYnzSjotVQyJjmtVg0bTcx6&listtoken=db6e707b6ece12b6c5f09776ef422b22&_ajax_nonce=17ca51e8a2`)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Set("Cookie", "wordpress_sec_1345aae3f89c108f9d0d100d1e068a48=loadhit.001%7C1532705947%7Cb0UnCTXDPfkGmsZaEJC2JGUapo3WaTMQAnYX4lATOul%7C3175c40ac3b8d8ae336429089149ea0d134e98c91d2c5d2f6c8f767b35190fb4; wordpress_cf_adm_use_adm=1; language=en; PHPSESSID=448jq3itik19v63khhqjole3sh; wordpress_test_cookie=WP+Cookie+check; wordpress_logged_in_1345aae3f89c108f9d0d100d1e068a48=loadhit.001%7C1532705947%7Cb0UnCTXDPfkGmsZaEJC2JGUapo3WaTMQAnYX4lATOul%7C9cb8a717345fce120c4b7e878bbfa4a91a42d848e6fc5b5ccba999829f8812a2; wp-settings-time-209=1532533210")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)
	log.Println("Payload : ", payload)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var drive driveResp
	if err := json.Unmarshal(body, &drive); err != nil {
		log.Println(err)
		return "", err
	}
	log.Println(drive.URL)

	////////////////////////////////////////////////////////
	const fileChunk = 5 * (1 << 20) // 1 MB, change this to your requirement

	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))

	log.Printf("Splitting to %d pieces.\n", totalPartsNum)

	var fileGurl driveFile
	for i := uint64(0); i < totalPartsNum; i++ {

		b3 := uint64(fileSize)

		b1 := uint64(i * uint64(fileChunk))
		b2 := uint64(0)

		if uint64(((i+1)*fileChunk)-1) > b3 {
			b2 = uint64(fileSize - 1)
		} else {
			b2 = uint64(((i + 1) * fileChunk) - 1)
		}
		b := "bytes " + strconv.FormatUint(b1, 10) + "-" + strconv.FormatUint(b2, 10) + "/" + strconv.FormatUint(b3, 10)

		// ------------------------------
		partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
		partBuffer := make([]byte, partSize)

		file.Read(partBuffer)

		// ------------------------------------
		putRequest, _ := http.NewRequest("PUT", drive.URL, strings.NewReader(string(partBuffer)))
		putRequest.Header.Add("Content-Range", b)
		putRequest.Header.Add("Content-Type", "application/zip")

		putResponse, err := http.DefaultClient.Do(putRequest)
		if err != nil {
			log.Fatalf("Unable to be post to Google API: %v", err)
			return "", err
		}

		defer putResponse.Body.Close()
		log.Println(putResponse.StatusCode, " => ", b)

		if putResponse.StatusCode == 200 {
			jsonResp, _ := ioutil.ReadAll(putResponse.Body)
			if err := json.Unmarshal(jsonResp, &fileGurl); err != nil {
				log.Println(err)
				return "", err
			}
		}
	}

	log.Println("G link : ", fileGurl.WebContentLink)
	return fileGurl.WebContentLink, nil
}
