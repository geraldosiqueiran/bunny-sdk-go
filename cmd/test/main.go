package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/geraldo/bunny-sdk-go/storage"
	"github.com/geraldo/bunny-sdk-go/stream"
)

const (
	// Global API Key (para Storage Zone Management e Stream Library Management)
	apiKey = "55f77bd6-6936-4370-9585-51c942028078b718e441-bf34-41dd-bc2d-252fe11f1b33"

	// Storage Zone credentials (para File Operations)
	storageZoneName     = "my-storage-test1"
	storageZonePassword = "65a1fec6-5090-49c7-8b67fff24236-9dca-47b4"
	storageZoneRegion   = storage.RegionSaoPaulo // br.storage.bunnycdn.com

	// Stream Library credentials (para Video Operations)
	streamLibraryID  = int64(588840)
	streamAPIKey     = "429f593d-87a6-48c6-a33966d930a5-549c-424d"
)

func main() {
	ctx := context.Background()

	fmt.Println("=== Bunny SDK Go - Teste de Rotas ===")
	fmt.Println()

	// Testar Storage Zone Management API
	testStorageAPI(ctx)

	// Testar Storage File Operations API
	testFileOperations(ctx)

	// Testar Stream API
	testStreamAPI(ctx)
}

func testStorageAPI(ctx context.Context) {
	fmt.Println("--- Storage API ---")

	client := storage.NewClient(apiKey)
	zones := client.Zones()

	// Listar Storage Zones
	fmt.Println("\n[1] Listando Storage Zones...")
	listResp, err := zones.List(ctx, &storage.ZoneListOptions{
		Page:    1,
		PerPage: 10,
	})
	if err != nil {
		log.Printf("Erro ao listar zones: %v\n", err)
	} else {
		fmt.Printf("Total de zones: %d\n", listResp.TotalItems)
		for _, zone := range listResp.Items {
			fmt.Printf("  - ID: %d, Nome: %s, Region: %s\n", zone.ID, zone.Name, zone.Region)
		}
	}

	fmt.Println()
}

func testFileOperations(ctx context.Context) {
	fmt.Println("--- Storage File Operations ---")

	fs := storage.NewFileService(storageZoneName, storageZonePassword, storageZoneRegion)

	// 1. Upload de arquivo
	fmt.Println("\n[1] Fazendo upload de arquivo de teste...")
	testContent := "Hello from Bunny SDK Go! - " + fmt.Sprintf("%d", 123456)
	err := fs.Upload(ctx, "sdk-test/test-file.txt", strings.NewReader(testContent), nil)
	if err != nil {
		log.Printf("Erro ao fazer upload: %v\n", err)
		return
	}
	fmt.Println("Upload realizado com sucesso!")

	// 2. Listar arquivos
	fmt.Println("\n[2] Listando arquivos no diretório 'sdk-test'...")
	files, err := fs.List(ctx, "sdk-test")
	if err != nil {
		log.Printf("Erro ao listar arquivos: %v\n", err)
	} else {
		fmt.Printf("Arquivos encontrados: %d\n", len(files))
		for _, f := range files {
			fmt.Printf("  - %s (%.2f KB, isDir: %v)\n", f.ObjectName, float64(f.Length)/1024, f.IsDirectory)
		}
	}

	// 3. Download de arquivo
	fmt.Println("\n[3] Fazendo download do arquivo...")
	reader, err := fs.Download(ctx, "sdk-test/test-file.txt")
	if err != nil {
		log.Printf("Erro ao fazer download: %v\n", err)
	} else {
		defer reader.Close()
		content, _ := io.ReadAll(reader)
		fmt.Printf("Conteúdo do arquivo: %s\n", string(content))
	}

	// 4. Deletar arquivo
	fmt.Println("\n[4] Deletando arquivo de teste...")
	err = fs.Delete(ctx, "sdk-test/test-file.txt")
	if err != nil {
		log.Printf("Erro ao deletar arquivo: %v\n", err)
	} else {
		fmt.Println("Arquivo deletado com sucesso!")
	}

	fmt.Println()
}

func testStreamAPI(ctx context.Context) {
	fmt.Println("--- Stream API ---")

	// Client com Global API Key para listar libraries
	client := stream.NewClient(apiKey)
	libraries := client.Libraries()

	// Listar Libraries
	fmt.Println("\n[1] Listando Video Libraries...")
	listResp, err := libraries.List(ctx, &stream.LibraryListOptions{
		Page:         1,
		ItemsPerPage: 10,
	})
	if err != nil {
		log.Printf("Erro ao listar libraries: %v\n", err)
	} else {
		fmt.Printf("Total de libraries: %d\n", listResp.TotalItems)
		for _, lib := range listResp.Items {
			fmt.Printf("  - ID: %d, Nome: %s\n", lib.LibraryID, lib.Name)
		}
	}

	// Testar operações de vídeo com Stream API Key específica da library
	testVideoOperations(ctx, streamLibraryID)

	fmt.Println()
}

func testVideoOperations(ctx context.Context, libraryID int64) {
	// Client com Stream API Key específica da library para operações de vídeo
	client := stream.NewClient(streamAPIKey)

	fmt.Printf("\n[2] Listando Vídeos da Library %d...\n", libraryID)

	videos := client.Videos(libraryID)
	videoList, err := videos.List(ctx, &stream.VideoListOptions{
		Page:         1,
		ItemsPerPage: 5,
	})
	if err != nil {
		log.Printf("Erro ao listar vídeos: %v\n", err)
		return
	}

	fmt.Printf("Total de vídeos: %d\n", videoList.TotalItems)
	for _, v := range videoList.Items {
		fmt.Printf("  - ID: %s, Título: %s, Estado: %s\n", v.VideoID, v.Title, v.State)
	}

	// Testar Collections
	fmt.Printf("\n[3] Listando Collections da Library %d...\n", libraryID)
	collections := client.Collections(libraryID)
	collList, err := collections.List(ctx, &stream.CollectionListOptions{
		Page:         1,
		ItemsPerPage: 5,
	})
	if err != nil {
		log.Printf("Erro ao listar collections: %v\n", err)
		return
	}

	fmt.Printf("Total de collections: %d\n", collList.TotalItems)
	for _, c := range collList.Items {
		fmt.Printf("  - GUID: %s, Nome: %s\n", c.GUID, c.Name)
	}
}

