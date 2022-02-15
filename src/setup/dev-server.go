package setup

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/tristanisham/colors"
)

func (b *DefaultBlog) Dev() error {
	fmt.Println(colors.As("Dev server started...", colors.LightBlue))
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					if err := b.buildArticles(); err != nil {
						log.Print(err)
					}
					log.Print(colors.As(event, colors.DarkGreen))

				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println(colors.As(err, colors.DarkRed))
			}
		}
	}()

	err = watcher.Add(os.Getenv("WRK_DIR") + "/src/posts")
	if err != nil {
		return err
	}
	if !<-done {
		return fmt.Errorf("watcher failed")
	}
	return nil
}
