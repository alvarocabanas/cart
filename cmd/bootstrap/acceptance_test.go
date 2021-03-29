package bootstrap_test

import (
	"bytes"
	"cart/cmd/bootstrap"
	"cart/internal/creator"
	"cart/internal/getter"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAcceptance(t *testing.T) {
	var (
		ctx         context.Context
		cancel      context.CancelFunc
		requireThat *require.Assertions
		server      *http.Server
		serverPort  = ":8889"
		client      *http.Client
		serverURL   string
	)

	setup := func() {
		var err error
		requireThat = require.New(t)
		ctx, cancel = context.WithCancel(context.Background())

		cfg := bootstrap.Config{ServerPort: serverPort}
		serverURL = "http://localhost" + serverPort

		server, err = bootstrap.InitializeServer(ctx, cfg)
		if err != nil {
			panic(err)
		}

		tr := &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		}
		client = &http.Client{Transport: tr}

		go func() {
			fmt.Printf("Launching Server")
			_ = server.ListenAndServe()

		}()
	}
	defer func() {
		_ = server.Shutdown(ctx)
		cancel()
	}()

	setup()

	t.Run("Given a server running and no items", func(t *testing.T) {
		t.Run("When we get the cart", func(t *testing.T) {
			resp, err := client.Get(serverURL)
			requireThat.NoError(err)
			defer resp.Body.Close()

			var dto getter.GetCartStatusDTO
			body, err := ioutil.ReadAll(resp.Body)
			requireThat.NoError(err)
			err = json.Unmarshal(body, &dto)
			t.Run("Then the response should be as expected", func(t *testing.T) {
				requireThat.NoError(err)
				requireThat.Equal(0, dto.TotalPrice)
			})
		})

		t.Run("When we add an item to the cart", func(t *testing.T) {
			itemDTO := creator.AddItemDTO{
				ItemID:   "book",
				Quantity: 8,
			}

			req, err := json.Marshal(itemDTO)
			requireThat.NoError(err)
			resp, err := client.Post(serverURL, "application/json", bytes.NewBuffer(req))
			requireThat.NoError(err)
			defer resp.Body.Close()

			t.Run("Then we get te correct response code", func(t *testing.T) {
				requireThat.NoError(err)
				requireThat.Equal(http.StatusCreated, resp.StatusCode)
			})
		})

		t.Run("When we get the cart", func(t *testing.T) {
			resp, err := client.Get(serverURL)
			requireThat.NoError(err)
			defer resp.Body.Close()

			var dto getter.GetCartStatusDTO
			body, err := ioutil.ReadAll(resp.Body)
			requireThat.NoError(err)
			err = json.Unmarshal(body, &dto)
			t.Run("Then the response should be as expected", func(t *testing.T) {
				requireThat.NoError(err)
				requireThat.Equal(480, dto.TotalPrice)
				requireThat.Equal(1, len(dto.Items))
			})
		})

		t.Run("When we add an item that doesn't exist", func(t *testing.T) {
			itemDTO := creator.AddItemDTO{
				ItemID:   "keyboard",
				Quantity: 8,
			}

			req, err := json.Marshal(itemDTO)
			requireThat.NoError(err)
			resp, err := client.Post(serverURL, "application/json", bytes.NewBuffer(req))
			requireThat.NoError(err)
			defer resp.Body.Close()

			t.Run("Then we get te correct response code", func(t *testing.T) {
				requireThat.NoError(err)
				requireThat.Equal(http.StatusNotFound, resp.StatusCode)
			})
		})

		t.Run("When we add a wrong amount", func(t *testing.T) {
			itemDTO := creator.AddItemDTO{
				ItemID:   "book",
				Quantity: -8,
			}

			req, err := json.Marshal(itemDTO)
			requireThat.NoError(err)
			resp, err := client.Post(serverURL, "application/json", bytes.NewBuffer(req))
			requireThat.NoError(err)
			defer resp.Body.Close()

			t.Run("Then we get te correct response code", func(t *testing.T) {
				requireThat.NoError(err)
				requireThat.Equal(http.StatusBadRequest, resp.StatusCode)
			})
		})
	})

	//Here all the cases should be contemplated, all the different response codes etc.
}
