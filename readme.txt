--- webp jpg processing timing ---
Compares the necessary processing times of jpg and webp for different scenarious.


Tests:
1. ImageMagick: Converting various file formats to jpg and webp
	1.1 jpg
	1.2 png
	1.3 webp
	1.4 tiff
	(formats can be added or removed in variable "m" in line 30)
2. libvips: Image manipulation of jpg and webp files from HD (1920px) to...
	2.1 resizing 500px
	2.2 resizing 860px
	2.3 resizing 1200px
	2.4 rotating
	2.5 force resizing
	(Tests can be added or removed at line 181ff)
	
Results:
Alle timing results are saved in the "results" folder as .txt-files.
The content of those .txt-files can be copied directly into the provided excel file into the corresponding sheet.
The graphical comparison is then done automatically.
	
Procedure:
The folder "originals" contains 40 random .jpg images (source: https://unsplash.com/). In the first step these images get formatted into jpg,png,webp,tiff images (formats can be added or removed in variable "m" in line 30). These are then saved into the folder "files_for_comparison". If they converted once, the variable "create_files" can be set to false. This saves time during the next runs.
The number of files that should be compared can be set with the variable "number_of_files_compared". The maximum is 40.
All results are the average of all images used. The more images used, the more stable the results will be, but the longher the runtime will be.
All image manipulation tests are being done for different image qualities. The image gets compressed at first and then all the tests are being done.
Some of the new files get saved in the folder "Zwischenspeicher", but get deleted again when they are not being used anymore.
