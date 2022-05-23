# CUDI (Clean Up Docker Images)

## Reduce Container Registry Storage

Container registries become large over time without cleanup. When a large number of images or tags are added:

* Fetching the list of available tags or images becomes slower.
* They take up a large amount of storage space on the server.

We recommend deleting unnecessary images and tags, and setting up a cleanup policy
to automatically manage your container registry usage.

## Cleanup policy

The cleanup policy is a scheduled job you can use to remove tags from the Container Registry.
For the project where it's defined, tags matching the regex pattern are removed.
The underlying layers and images remain.

## How the cleanup policy works
The cleanup policy collects all tags in the Container Registry and excludes tags
until only the tags to be deleted remain.

The cleanup policy searches for images based on the tag name. Support for the full path has not yet been implemented, but would allow you to clean up dynamically-named tags.

The cleanup policy:

1. Collects all matching tags from **removeTags** list for a given repository in a list.
2. Excludes from the list any tags matching the **keepTags** value (tags to preserve).
3. Orders the remaining tags by created_date.
4. Excludes from the list the N tags based on the **keepNumber** value (Number of tags to retain).
5. Excludes from the list the tags more recent than the **olderThan** value (Expiration interval).
6. Finally, the remaining tags in the list are deleted from the Container Registry.

Valid values for type:

* together
* separately

Valid units for olderThan:

* y: (year)
* m: (month)
* w: (week)
* d: (day)
* h: (hour)
* min: (minute)
* s: (second)