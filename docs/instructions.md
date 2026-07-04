# Instructions

Included in the zip is a `commits.csv` file of commits into the default branches of GitHub repositories from teams working in an inner source model at Flutter. This is anonymised real data with the fields:

| Field        | Description                                                  |
| ------------ | ------------------------------------------------------------ |
| `timestamp`  | The unix timestamp of the commit                             |
| `user`       | The GitHub username of the commit author. This data is unreliable and blank indicates the author's identity is unknown. |
| `repository` | The repository name which the commit was pushed to.          |
| `files`      | The number of files changes by the commit                    |
| `additions`  | The number of line additions in this commit                  |
| `deletions`  | The number of deletions in this commit                       |

The file contains 100 days of commits. 

## Your Task

Using a programming language of your own choice:

1. Design and implement an algorithm that will give each repository an activity score.

2. Document this algorithm in a markdown file, and any directions required to run your implementation of it.

3. Use this algorithm to rank the repositories and include the top 10 "most active" repositories by your definition in the documentation.

4. Zip up the documentation and your implementation and return it to complete your test. 

The task should take no longer than 1-2 hours. Note that you are welcome to open and explore the csv data in Excel if you wish, but your ranking algorithm should be implemented in a programming language for automated execution.