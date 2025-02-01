# Auto bid generator
Generates bids with AI for freelancer.com projects.

## Instalation
make executable
`go build`

export openai api key
`export AI_APIKEY=api-key`

## Usage
`./bidgen -u freelancer-project-url`

Example:
`./bidgen -u https://www.freelancer.com.ar/projects/iphone-app-development/commerce-Mobile-App-Development-39023417/details`

## Prompt
if you want to change the prompt, simply modify the `const prompt ` at the top of `prompt.go` file.

# TO DO
- [x] temperature flag
- [ ] model flag 
- [ ] -x flag to copy directly to the clipboard
- [ ] add the project title and more information from the freelancer api
- [ ] test with openai models (they decline my card)
- [ ] unit test for bid creation
- [ ] stream of tokens? or some load animation while waiting for the bid
- [ ] better prompt
- [ ] test with english projects

# IDEAS
- extension for the web browser
