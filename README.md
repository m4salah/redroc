# Redroc

Scalabale image service that allow users download, search, and upload images.

## Services

We have two type of services gRPC, and RESTful services.

- ### gRPC services

    We cannot call those services directly, those services required authentication from the caller in Google Cloud (Cloud Run).

1. #### Download

    Download service that allow user to get download the image by image name.

    the service proto is

    ```proto
        message DownloadPhotoRequest {
            string img_name = 1;
        }

        message DownloadPhotoResponse {
            bytes img_blob = 1;
        }

        service DownloadPhoto {
            /*
            * RPC for download a photo
            */
            rpc Download(DownloadPhotoRequest) returns (DownloadPhotoResponse);
        }
    ```

2. #### Upload

    Upload service that allow user upload image with it's metadata to search for image later.

    Require username, and file.

    Optional hashtags.

    the service proto is

    ```proto
        message UploadImageRequest {
            string obj_name = 1;
            bytes image = 2;
        }

        message UploadImageResponse {}

        message CreateMetadataRequest {
            string obj_name = 1;
            string user = 2;
            repeated string hashtags = 3;
        }
        message CreateMetadataResponse {}

        service UploadPhoto {
            /*
            * RPC for upload a photo to the image database
            */
            rpc Upload(UploadImageRequest) returns (UploadImageResponse);

            /*
            * RPC for create hashtag-image mapping in the metadata database
            */
            rpc CreateMetadata(CreateMetadataRequest) returns (CreateMetadataResponse);
        }

    ```

3. #### Search

    Search service that allow user to search for images with hashtags attached with the image.

    if search keyword is empty it will return latest image uploaded.

    the service proto is

    ```proto
        message GetThumbnailImagesRequest {
        // if keyword=="latest", return recent photo
        // in the service, we will update metadata such as download_times accordingly
        string search_keyword = 1;
        }

        message GetThumbnailImagesResponse {
        // get the storage image-serving address and return
        repeated string storage_url = 1;
        }

        service GetThumbnail {
        /*
            RPC for getting the UIDs of images relevant to the keyword
        */
        rpc GetThumbnail(GetThumbnailImagesRequest)
            returns (GetThumbnailImagesResponse);
        }
    ```
