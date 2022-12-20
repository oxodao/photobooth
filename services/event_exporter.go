package services

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/oxodao/photobooth/models"
	"github.com/oxodao/photobooth/orm"
	"github.com/oxodao/photobooth/utils"
)

type EventExporter struct {
	event *models.Event
}

func NewEventExporter(event *models.Event) EventExporter {
	return EventExporter{event}
}

func (ee EventExporter) setEventExporting(exp bool) error {
	ee.event.Exporting = exp
	err := orm.GET.Events.Save(ee.event)
	if err != nil {
		fmt.Println("Failed to set the exporting state")
		return err
	}

	return nil
}

func (ee EventExporter) Export() (*models.ExportedEvent, error) {
	if ee.event.Exporting {
		return nil, errors.New("can't export an event that is already exporting")
	}

	if err := ee.setEventExporting(true); err != nil {
		return nil, err
	}

	exportTime := time.Now()

	basepath := fmt.Sprintf("images/%v/exports/", ee.event.Id)
	err := utils.MakeOrCreateFolder(basepath)
	if err != nil {
		if err2 := ee.setEventExporting(true); err2 != nil {
			return nil, err
		}

		return nil, err
	}

	basefilepath := exportTime.Format("20060201-150405") + ".zip"
	filepath := utils.GetPath(basepath + basefilepath)

	archive, err := os.Create(filepath)
	if err != nil {
		if err2 := ee.setEventExporting(true); err2 != nil {
			return nil, err
		}

		return nil, err
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	//#region Exporting non-unattended images
	images, err := orm.GET.Events.GetImages(ee.event, false)
	if err != nil {
		if err2 := ee.setEventExporting(true); err2 != nil {
			return nil, err
		}

		return nil, err
	}

	for _, i := range images {
		imagePath := utils.GetPath(fmt.Sprintf("images/%v/pictures/%v.jpg", ee.event.Id, i.Id))
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			fmt.Printf("Failed to locate image %v from event %v\n", i.Id, ee.event.Id)
			continue
		}

		fr, err := os.Open(imagePath)
		if err != nil {
			fmt.Printf("Failed to open the image %v for the event %v: %v\n", i.Id, ee.event.Id, err)
			continue
		}

		fw, err := zipWriter.Create((time.Time(i.Date)).Format("20060201-150405") + ".jpg")
		if err != nil {
			fmt.Printf("Failed to create the image %v for the event %v in the zip file: %v\n", i.Id, ee.event.Id, err)
			continue
		}

		if _, err := io.Copy(fw, fr); err != nil {
			fmt.Printf("Failed to copy the image %v for the event %v in the zip file: %v\n", i.Id, ee.event.Id, err)
			continue
		}

		fr.Close()
	}
	//#endregion

	//#region Exporting unattended images
	unattendedRoot := utils.GetPath(fmt.Sprintf("images/%v/unattended/", ee.event.Id))
	outvid := unattendedRoot + "/000_recap.mp4"

	if _, err := os.Stat(outvid); !os.IsNotExist(err) {
		os.Remove(outvid)
	}

	cmd := exec.Command("bash", "-c", "ffmpeg -framerate 6 -pattern_type glob -i '*.jpg' -c:v libx264 -vf \"scale=iw*min(1080/iw\\,1080/ih):ih*min(1080/iw\\,1080/ih), pad=1080:1080:(1080-iw*min(1080/iw\\,1080/ih))/2:(1080-ih*min(1080/iw\\,1080/ih))/2,fps=30000/1001,format=yuv420p\" 000_recap.mp4")
	cmd.Dir = unattendedRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		if err2 := ee.setEventExporting(false); err2 != nil {
			return nil, err
		}
		return nil, err
	}

	fr, err := os.Open(outvid)
	if err != nil {
		fmt.Printf("Failed to open the recap video for the event %v: %v\n", ee.event.Id, err)
		if err2 := ee.setEventExporting(false); err2 != nil {
			return nil, err
		}

		return nil, err
	}

	fw, err := zipWriter.Create("000_recap.mp4")
	if err != nil {
		fmt.Printf("Failed to create the recap video for the event %v in the zip file: %v\n", ee.event.Id, err)
		if err2 := ee.setEventExporting(false); err2 != nil {
			return nil, err
		}
		return nil, err
	}

	if _, err := io.Copy(fw, fr); err != nil {
		fmt.Printf("Failed to copy the recap video for the event %v in the zip file: %v\n", ee.event.Id, err)
		if err2 := ee.setEventExporting(false); err2 != nil {
			return nil, err
		}
		return nil, err
	}

	fr.Close()

	if _, err := os.Stat(outvid); !os.IsNotExist(err) {
		os.Remove(outvid)
	}
	//#region

	//#region Adding a json file with some data about the event
	data := map[string]interface{}{
		"id":                    ee.event.Id,
		"name":                  ee.event.Name,
		"author":                ee.event.Author,
		"date":                  ee.event.Date,
		"location":              ee.event.Location,
		"amt_images_handtaken":  ee.event.AmtImagesHandtaken,
		"amt_images_unattended": ee.event.AmtImagesUnattended,
	}
	jsonData, _ := json.MarshalIndent(data, "", "  ")

	fw, err = zipWriter.Create("001_infos.json")
	if err != nil {
		fmt.Println("Failed to add the info json")
		if err2 := ee.setEventExporting(false); err2 != nil {
			return nil, err
		}
		return nil, err
	}

	if _, err := io.Copy(fw, bytes.NewReader(jsonData)); err != nil {
		fmt.Printf("Failed to copy the info json for the event %v in the zip file: %v\n", ee.event.Id, err)
		if err2 := ee.setEventExporting(false); err2 != nil {
			return nil, err
		}
		return nil, err
	}
	//#endregion

	exportTimestamp := models.Timestamp(exportTime)
	ee.event.LastExport = &exportTimestamp
	err = ee.setEventExporting(false)
	if err != nil {
		return nil, err
	}

	// Insert the built
	exportedEvent, err := orm.GET.Events.InsertExportedEvent(ee.event, basefilepath)
	if err != nil {
		return nil, err
	}

	return exportedEvent, nil
}
