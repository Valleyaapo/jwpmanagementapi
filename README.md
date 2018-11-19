# JWP Management API golang URL creation
A golang library you can use to create and call valid JWP management API requests.</br>
Use your own API_SECRET and API_KEY as environment variables</br>
Actual API documentation: </br>
https://developer.jwplayer.com/jw-platform/reference/v1/methods/accounts/tags/index.html </br>
## Usage
TagManager takes name of the tag as a slice of strings and name of the method as string. `Create` and `Delete` take only one parameter and return nil. `Update` takes two parameters. First is the old name to be replaced and second is the new name and returns nil. `List` takes no parameters, but returns a slice of strings of all tags associated with the account. List is sorted alphabetically.

## Disclaimer
I have no connection to Longtail Ad Solutions or JWPlayer. This is purely a personal project.
