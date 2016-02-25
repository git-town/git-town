# Script to build and run the docker instance
#provide it the absolute path of the folder you want to save the deb file to
echo "Please enter in the directory that you want to save the deb file to: "
read -r user_dir
echo "Saving to directory: $user_dir"

docker build -t allonsy/gittown_build:latest .
docker run -v "$user_dir":/shared allonsy/gittown_build
