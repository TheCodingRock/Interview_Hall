/*
Author: Stoney Alexander Hall // @TheCodingRock
Created: 7/26/22
Purpose: Take in large list of date/time data. Remove any bad data. Then write the clean data to a new document
*/
package main

import (
	//For scanning
	"bufio"
	//For string formating
	"fmt"
	//For file creation
	"os"
)

// Asks user for their file and sends adress to the read_Edit_Write_File function.
func main() {
	fmt.Println("Please enter the [fileName.Extension] of where your data is located: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		panic(err)
	}
	fmt.Println("Please enter the [fileName.Extension] of the document you would like to create:")
	scanner2 := bufio.NewScanner(os.Stdin)
	scanner2.Scan()
	errr := scanner2.Err()
	if err != nil {
		panic(errr)
	}
	readEditWriteFile(scanner.Text(), scanner2.Text())
}

/*
Reads in data from file provided by the user.
Sends data to be scrubbed and edited by validation function.
Writes clean data to a new document.
*/
func readEditWriteFile(filePath string, destinationFile string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanWords)

	var dateTimes []string
	pointer := &dateTimes

	for fileScanner.Scan() {
		dateTimes = append(dateTimes, fileScanner.Text())
	}
	if err := fileScanner.Err(); err != nil {
		panic(err)
	}

	validate(pointer)
	checkForDuplicates(pointer)
	writeFile(pointer, destinationFile)
	fmt.Println("Process Complete! Your results can be found in a file called " + destinationFile + ".")
}

//Identifies and removes duplicate entries.
func checkForDuplicates(original *[]string) {
	temp := *original
	instance := map[string]bool{}
	result := []string{}
	for element := range temp {
		if instance[temp[element]] != true {
			instance[temp[element]] = true
			result = append(result, temp[element])
		}
	}
	*original = result
}

//Creates a new .txt file and writes finished data to it
func writeFile(cleanData *[]string, name string) {
	file, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var temp []string = *cleanData
	var holding string
	for i := 0; i < len(temp); i++ {
		holding = temp[i]
		if i == (len(temp) - 1) {
			_, err = file.WriteString(holding)
			if err != nil {
				panic(err)
			}
		} else {
			_, err = file.WriteString(holding + "\n")
			if err != nil {
				panic(err)
			}
		}
	}
}

/*
Steps through imported data and validates it against the format YYYY-MM-DDThh:mm:ss['Z'][+hh:mm][-hh:mm]
using each characters ascii code for type validation.
Bad data indexes are sent to the deletion function which then removes them from the list.
*/
func validate(original *[]string) {
	var temp []string = *original
	var holding rune
	x := 0
	for i := 0; i < len(temp); i++ {
		tempHold := temp[i]
		holding = 0
		// Item has to be either 20 characters long using 'Z' or 25 long using '+-hh:mm'
		if !(len(tempHold) == 20 || len(tempHold) == 25) {
			delEntry(i, &temp)
			x = 1
			i--
		} else {
			for j, char := range tempHold {
				switch j {
				// YYYY
				case 0, 1, 2, 3:
					// Has to be integer 0-9
					if !(char >= 48 && char <= 57) {
						delEntry(i, &temp)
						x = 1
					}
				// '-'
				case 4:
					if char != 45 {
						delEntry(i, &temp)
						x = 1
					}
				// MM
				case 5, 6:
					// Has to be integer 0-1
					if j == 5 {
						if !(char == 48 || char == 49) {
							delEntry(i, &temp)
							x = 1
						} else {
							holding = char
						}
						// Has to be integer 0-9
					} else if j == 6 {
						if !(char >= 48 && char <= 57) {
							delEntry(i, &temp)
							x = 1
							// Can't be more than 12
						} else if holding == 49 {
							if !(char >= 48 && char <= 50) {
								delEntry(i, &temp)
								x = 1
								// Can't be 00
							} else if holding == 48 {
								if char == 48 {
									delEntry(i, &temp)
									x = 1
								}
							}
						}
					}
				// '-'
				case 7:
					if char != 45 {
						delEntry(i, &temp)
						x = 1
					}
				// DD
				case 8, 9:
					// Has to be integer 0-3
					if j == 8 {
						if !(char >= 48 && char <= 51) {
							delEntry(i, &temp)
							x = 1
						} else {
							holding = char
						}
						// Has to be integer 0-9
					} else if j == 9 {
						if !(char >= 48 && char <= 57) {
							delEntry(i, &temp)
							x = 1
							// Can't be more than 31
						} else if holding == 51 {
							if !(char == 48 || char == 49) {
								delEntry(i, &temp)
								x = 1
							}
							// Can't be 00
						} else if char == 48 {
							if holding == 48 {
								delEntry(i, &temp)
								x = 1
							}
						}
					}
				// 'T'
				case 10:
					if char != 84 {
						delEntry(i, &temp)
						x = 1
					}
				// hh
				case 11, 12:
					// Has to be an integer 0-2
					if j == 11 {
						if !(char >= 48 && char <= 50) {
							delEntry(i, &temp)
							x = 1
						} else {
							holding = char
						}
						// Has to be integer 0-9
					} else if j == 12 {
						if !(char >= 48 && char <= 57) {
							delEntry(i, &temp)
							x = 1
							// Can't be more than 23
						} else if holding == 50 {
							if !(char >= 48 && char <= 51) {
								delEntry(i, &temp)
								x = 1
							}
						}
					}
				// ':'
				case 13:
					if char != 58 {
						delEntry(i, &temp)
						x = 1
					}
				// mm
				case 14, 15:
					// Has to be integer 0-5
					if j == 14 {
						if !(char >= 48 && char <= 53) {
							delEntry(i, &temp)
							x = 1
						}
						// Has to be integer 0-9
					} else if j == 15 {
						if !(char >= 48 && char <= 57) {
							delEntry(i, &temp)
							x = 1
						}
					}
				// ':'
				case 16:
					if char != 58 {
						delEntry(i, &temp)
						x = 1
					}
				// ss
				case 17, 18:
					// Has to be integer 0-5
					if j == 17 {
						if !(char >= 48 && char <= 53) {
							delEntry(i, &temp)
							x = 1
						}
						// Has to be integer 0-9
					} else if j == 18 {
						if !(char >= 48 && char <= 57) {
							delEntry(i, &temp)
							x = 1
						}
					}
				// 'Z' or '+' or '-'
				case 19:
					if !(char == 90 || char == 43 || char == 45) {
						delEntry(i, &temp)
						x = 1
					}
				// hh
				case 20, 21:
					// Has to be integer 0-2
					if j == 20 {
						if !(char >= 48 && char <= 50) {
							delEntry(i, &temp)
							x = 1
						} else {
							holding = char
						}
						// Has to be integer 0-9
					} else if j == 21 {
						if !(char >= 48 && char <= 57) {
							delEntry(i, &temp)
							x = 1
							// Can't be more than 23
						} else if holding == 50 {
							if !(char >= 48 && char <= 51) {
								delEntry(i, &temp)
								x = 1
							}
						}
					}
				// ':'
				case 22:
					if char != 58 {
						delEntry(i, &temp)
						x = 1
					}
				// mm
				case 23, 24:
					// Has to be integer 0-5
					if j == 23 {
						if !(char >= 48 && char <= 53) {
							delEntry(i, &temp)
							x = 1
						}
						// Has to be integer 0-9
					} else if j == 24 {
						if !(char >= 48 && char <= 57) {
							delEntry(i, &temp)
							x = 1
						}
					}
				default:
					fmt.Println("WARNING: Something went wrong during processing. Results may not be accurate!")
				}
				// Accounts for missing indnex when item is deleted
				if x == 1 {
					x = 0
					i--
					break
				}
			}
		}
	}
	// Updates the main list
	*original = temp
}

/*
Recieves a pointer for the list of elements needing editing as well as the index of the target for removal.
Removes element at index. Reslices list to acomidate for lost index.
*/
func delEntry(i int, original *[]string) {
	var temp []string = *original
	if i < 0 || i > len(temp) {
		fmt.Println("WARNING: Something went wrong during processing. Results may not be accurate!")
	} else {
		newLength := 0
		for index := range temp {
			if i != index {
				temp[newLength] = temp[index]
				newLength++
			}
		}
		*original = temp[:newLength]
	}
}
