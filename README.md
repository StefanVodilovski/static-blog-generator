# static-blog-generator
A CLI application for static blog generation

The application is made in multiple itterations and each version is available on the different branches:
  * main -> the final version of the app
  * minimal-solution -> inital version that setisfies the minimal solution requirements
  * stretch-goal-pagination -> added pagination to the minimal solution.
  * stretch-goal-pagination-sort -> added sorting to the blog generator.
  * dockerized-app -> added docker to containerize the application and make it work on all OS (merged with main)


# Installation
The system has been containerized in order to support different OS environments.  To be able to use the CLI application you need to clone this repository and build the docker image using the following command: 

```
docker build . -t gen-blog 
```

# Usage
The application runs in a containerized environment. 
To use the CLI for generating HTML files from .MD files run the following command:

```
docker run -v <your_input_folder>:/app/input -v <your_output_folder>:/app/output gen-blog generate --input /app/input --output /app/output --title <your_title> --posts-per-page <optional_posts_per_page>
```

where <your_input_folder> contains .md files and <your_output_folder>, which must be an emtpy folder, will contain the output HTML files.

