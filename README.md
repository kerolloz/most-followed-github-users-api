<h1 align="center">Most Followed GitHub Users (API)</h1>
<p align="center">
  Returns a list of the top 50 most followed GitHub users in a particular country
</p>
<p align="center">
  <img src="https://t.ly/plOs" alt="go" />
  <img src="https://t.ly/MMyE" alt="deta" />
  <img src="https://t.ly/ljN-f" alt="github" />
</p>
<p align="center">
  <img src="https://user-images.githubusercontent.com/36763164/232241911-c0771365-225c-49fc-b22a-8351b54336b9.png" width="500px" />
</p>

## Getting Started

To get started with the project, follow these steps:

1. Clone the repository  
1. Navigate to the project directory
1. Set the [required environment variables](#environment-variables)
1. Install the project dependencies using `go get ./...`  
1. Build the project using `go build -o main ./app`
1. Run the project using `./main`

## Environment Variables

The project requires certain environment variables to be set in order to function properly. You can find a list of required variables in the `.env.example` file.

The project is flexible when it comes to loading the environment variables.
You can use any of the following ways to set the environment variables:

- Create a `.env` file in the root directory of the project and add your variables there. The project will automatically read these variables on startup.
- Set environment variables directly in your shell using the `export` command. For example, you could run `export GITHUB_API_TOKEN=your_token` to set the `GITHUB_API_TOKEN` variable. Note that this method will only set the variable for the current shell session.
- When running the compiled executable file, you can set environment variables inline like this: `GITHUB_API_TOKEN=your_token ./main`. This will set the `GITHUB_API_TOKEN` variable specifically for the execution of that command.

## Issues

If you encounter any issues with the project, please report them on the GitHub issue tracker.

## Contributing

If you would like to contribute to the project, please follow these steps:

1. Fork the repository
2. Create a new branch for your feature or bug fix
3. Make your changes and commit them
4. Push your changes to your forked repository
5. Open a pull request
