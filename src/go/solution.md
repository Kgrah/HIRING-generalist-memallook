**MEMALLOOK**

Instructions:
1. Clone repo
2. Make sure go is installed: https://golang.org/doc/install
3. navigate to src/go in your terminal
4. Build the go with go build main.go
5. invoke the app directly through the main binary --ex: `main show`
**Note** Due to time restrictions, this app is not fully functional. However I do believe the algorithmic code in my memallook package expresses the idea of 
the memory manager itself.

References:
1. Go docs: https://golang.org/doc/ -- for documentation surrounding specific golang apis
2. A coworker named Chris to discuss general strategies around memory allocation (naive hunting algorithms vs smarter mem allocation algorithms), data structures
to express memory allocation, and methods for persistence.

Design Choices:
1. Persisting the data as a JSON.
  a. This requires reading in the "state" of the program from a file on disk every time a command is invoked, rather than holding a process open with 
      in-memory data allocations
  b. This requires re-writing the entire JSON object representing the memory buffer on *every* state change. This would not work with concurrency.
  c. This was the fastest solution given personal time constraints

2. Naive searching algorithm for finding free space in the buffer
  a. Intended to combine with compaction when no free space is found, but ran out of time.
  b. Lended itself well to an easy to think about compaction algorithm that involved simply shifting memory tags forward up the memory slice.
  c. Was slow in most cases (minimum o(n) n being number of bytes currently being allocated)
  
3. Leaving allocated blocks set to their value on deallocation rather than zeroing out values.
  a. Improved runtime complexity -- we do not have to scan the slice for data potentially spread across multiple pages to be zeroed out, but instead can just
      delete the tag describing the memory allocation from the tag map.
  b. In production systems this could be a security vulnerably if an attacker was given access to the memory buffer.

Libraries and frameworks:
<br></br>
I did not pull in any libraries or frameworks outside of what Golang offers out of the box.
