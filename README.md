# gobase
Go Container Template + VSCode Dev Container Config

## What is this?
This is a utility repository with the most bare-bones setup to begin developing with containers for Go.

## Getting Started

Gobase contains configurations for **development containers with VSCode**, so you can just get into developing without worrying about container configuration too much.

1. 🐳 Make sure you have [Docker](https://docs.docker.com/get-docker) installed on your system.

    * 👉 Verify this by running:
        ```bash
        $ docker version
        ```

   🐳 Make sure you also have [Docker Compose](https://docs.docker.com/compose/install/).

    * 👉 Verify this by running:
        ```bash
        $ docker-compose version
        ```

    <br>
2. 📥 Install the [Remote - Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) extension in your VSCode instance.

    Verify that it works by clicking on the whale icon in the left sidebar. 
    
    No errors should show up.
    
    <br>

3. 📥 Clone this repository somewhere where you want to initialize your new project.

    ```bash
    # --depth=1 means to only take commits from 1 version back (shallow copy)
    $ git clone --depth=1 git@github.com:taskcollect/gobase.git your_project_folder
    # go into the directory
    $ cd your_project_folder
    ```
    
    OR **YOU CAN USE GITHUB'S TEMPLATE FEATURE!**

4. 🔗 Remove the links to this repository. For this, you can just remove the .git and start anew:

    ```bash
    # remove links to gobase
    $ rm -rf .git/
    # make a new repo in this folder
    $ git init
    # add your remote
    $ git remote add origin git@github.com:taskcollect/your_go_repo
    ```
5. 🌐 Change `.devcontainer/docker-compose.yml` to add potential services you might need running in order for the container to start.

    * 🐘 An example of a simple Postgres database server is included. You can try uncommenting the service definition and adding it to the "depends on" section of the devcontainer.

6. 📦 Change `.devcontainer/devcontainer/Dockerfile` to run extra commands you might want, such as extra system packages or go dependencies.

7. ✅ Run the devcontainer! 

    1. <kbd>Ctrl/Cmd</kbd> + <kbd>Shift</kbd> + <kbd>P</kbd>

    2. Type **`Reopen Folder In Container`**, hit <kbd>Enter</kbd>.

    3. VSCode should now be building your container!

    4. Once it's done, you'll be dropped into a shell in the container as the `vscode` user. You also have access to Git in the container!

8. 💻 Develop to your heart's content. You should probably use `go run` for this.

    For example, you can run the included simple webserver like this:

    ```bash
    # in the container at ~/workspace
    $ cd ./src
    $ go run .
    ```

    And you should be able to view the website hosted by it at [localhost:2000](https://localhost:2000), with VSCode automatically detecting the port and forwarding it to your host.

9. 🚫 To exit the container, open the command picker (the same way as in step 7) and
    instead type **Reopen Folder Locally**.

9.  👍 Once you're finished with developing, you may want to test how your container
    would run in production, and set up a Docker image for this.

    ---

    👉 The production-grade multistage dockerfile is the Dockerfile in the root of the entire project. Don't confuse it with the development Dockerfile in .devcontainer!

    ---

    You should prepare it with any dependencies it might need in the build stage, and play around with it until it's stable.

    Once you've tested it and you're sure it works, build and tag it for production use.

    ```bash
    # replace your-image:version with something else.
    # you can also just omit everything past : to tag as latest
    
    # replace <edition> with either alpine or distroless
    $ docker build . -f Dockerfile.<edition> -t your-image:version
    ```

    For example, this is how you'd tag & test the example web server for production.
    ```bash
    # build the container image "gobase-example" from the distroless variant
    $ docker build . -f Dockerfile.distroless -t gobase-example
    # run, binding port 9000 on host to 2000 in the container
    $ docker run --rm -p 9000:2000 gobase-example
    # the app should be accessible on the host, at http://localhost:9000
    ```

    You can alternatively run `make distroless` or `make alpine`. You should check out the [Makefile](./Makefile) to see what that actually does. Remember to change the image name! 


That's it!

Note: You'll probably want to delete this file as well before pushing changes to another git repository. Happy hacking!
