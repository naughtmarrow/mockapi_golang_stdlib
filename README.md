# MOCK API

The overall purpose of this API is to provide the backend capabilities for a game development site that includes the usage of dev logs written with Markdown to give news and updates about upcoming games and current development projects.

It is built with mostly Golang's standard library, but it includes some external libraries found in the go.mod file.

The API's main endpoints are typical CRUD operations with some extra features.
Most send and receive JSON in some form, with the given tags.
The API's endpoints are as follows:

- Blogs:
    - POST → ^/blogs/*$
        - Takes in a blog object with the following JSON tags:
        "id"
        "title"
        "description"
        "upload_date"
        "update_date"
        "link_to_md"
        "link_to_jsx"
        - Returns the JSON object back

    - GET → ^/blogs/*$
        - Returns a list of all blogs in JSON

    - GET → ^/blogs/([0-9]+)
        - Returns the blog with the given ID in JSON

    - GET → ^/blogs/name/([a-zA-Z0-9_.-]*)$
        - Returns the blog with the given name in JSON

    - DELETE → ^/blogs/([0-9]+)
        - Delete the blog with the given ID in JSON

    - POST → ^/blogs/file/([0-9]+)
        - Receives the markdown file and saves it, it also parses the file and creates an HTML counterpart
        - Takes in the file itself as a multipart-form
        - Returns an HTML divisor with the parsed Markdown

    - GET → ^/blogs/file/([0-9]+)
        - Returns the HTML form of the Markdown file saved on the system

    - POST → ^/blogs/admin/fileform/*$
        - Takes in the ID of the object and returns HTML with a form to upload a file to the endpoint with that ID
        - Used in conjunction with the POST file upload API endpoint to get the correct ID

- Tags:
    - POST → ^/tags/*$
        - Takes in a tag object with the following JSON tags:
        "id"
        "title"
        "color"
        "link_to_svg"
        "blog_list" -- OPTIONAL
        - Returns the JSON object back

    - GET → ^/tags/*$
        - Returns a list of all tags in JSON

    - GET → ^/tags/([0-9]+)
        - Returns the tag with the given ID in JSON

    - GET → ^/tags/name/([a-zA-Z0-9_.-]*)$
        - Returns the tag with the given name in JSON

    - DELETE → ^/tags/([0-9]+)
        - Delete the tag with the given ID in JSON

    - GET → ^/tags/bloglist/([0-9]+)
        - Returns the blog list of the tag with the given ID in JSON 

    - POST → ^/tags/bloglist/([0-9]+)
        - Adds the blog with the ID sent in the body of the request to the list of the tag with the ID given in URL takes in the following JSON tags:
        "blog_id"

    - DELETE → ^/tags/bloglist/([0-9]+)
        - Deletes the blog with the ID sent in the body of the request from the list of the tag with the ID given in URL takes in the following JSON tags: 
        "blog_id"

- Users:
    - POST → ^/users/*$
        - Takes in a users object with the following JSON tags:
        "id"
        "username"
        "password"

    - GET → ^/users/*$
        - Returns a list of all users in JSON

    - GET → ^/users/([0-9]+)
        - Returns the users with the given ID in JSON

    - GET → ^/users/name/([a-zA-Z0-9_.-]*)$
        - Returns the users with the given name in JSON

    - PUT → ^/users/([0-9]+)
        - Takes in a key and a value in JSON form with the following JSON tags:
        "key"
        "value"
        - Key must be either "password" or "username" depending on what you want to update.
        - Value is the new value to update the key with.

    - DELETE → ^/users/([0-9]+)
        - Delete the users with the given ID in JSON

### TO DO:
- Add PUT methods to all the keys in tags and blog
