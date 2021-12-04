package main

import (
	"archive/tar"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

// Largely inspired by https://golangdocs.com/tar-gzip-in-golang

func Tar(output string, root string, inputs []string) error {
	// Create the output tar file
	tarfile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	// Create a writer for the tar file
	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	for _, input := range inputs {
		err := addToTarball(input, root, tarball)
		if err != nil {
			return err
		}
	}

	return nil
}

func addToTarball(source string, root string, tarball *tar.Writer) error {
	// Open the file to be added to written out to the tarball
	file, err := os.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()
	// Get metadata about the file
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	// Create a tar header using the metadata
	header, err := tar.FileInfoHeader(stat, stat.Name())
	if err != nil {
		return err
	}

	// // Internet says this is needed to properly preserve the directory structure
	header.Name = path.Join(root, stat.Name())
	if err := tarball.WriteHeader(header); err != nil {
		return err
	}

	if stat.IsDir() {
		return nil
	}
	_, err = io.Copy(tarball, file)
	return err
}

// Creates a "flat" tar containing a list of specified files stored under the specified root dir.
func main() {
	var output string
	flag.StringVar(&output, "output", "output.tar", "Output filename")
	var root string
	flag.StringVar(&root, "root", "", "Directory in which to store files inside the tarball")
	flag.Parse()
	inputs := flag.Args()

	if len(inputs) == 0 {
		fmt.Println("No inputs were specified, bailing out.")
		os.Exit(1)
	}

	fmt.Printf("Creating a tarball %s from:\n%s\n", output, strings.Join(inputs, "\n"))

	err := Tar(output, root, inputs)

	if err != nil {
		log.Fatalln("Error creating archive:", err)
	}

	fmt.Println("Done")
}
