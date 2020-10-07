Requirement
-----------

Allow cold storage (config) files replication on a cluster of two nodes

Architecture
------------

    |    Node 1            |         |  Node 2              |
    |   __________         |         |         __________   |
    |  | MGMT API |        |         |        | MGMT API |  |
    |  |__________|_____   |         |   _____|__________|  |
    |  | entangled|     |  |         |  |     | entangled|  |
    |  |          | M2M |  |         |  | M2M |          |  |
    |  |          | API |<-|---------|->| API |          |  |
    |  |__________|_____|  |         |  |_____|__________|  |
    |   ____|_____         |         |         ____|_____   |
    |  | inotify  |        |         |        | inotify  |  |
    |  |__________|        |         |        |__________|  |

Constraints/Assumptions
-----------------------

1. The tool only syncs single files (not whole folders). 
It is designed specifically in this way to avoid misusage (e.g. RDBD replacing).
2. If a file is listed to be monitored, when entangled is started it requires the path to that file to exist already.
3. If the file to be monitored does not exist, a warning is emitted, and if the other peer requires alignment, the file is created with 644 permissions
4. The tool does not resolve a split-brain

Use Cases
---------

#### Bootstrap

When entangled is started, it registers all files that need to be monitored and calculates their md5 checksum for future changes recognition.
If the two nodes files are misaligned, the misalignment is preserved. In other words, the tool does not resolve a split-brain.
Split-brain resolution is out of scope and delegated to manual procedures.

#### FS file change

When a FS monitored file is changed, if the new md5 checksum differs from cached one, the sync procedure is issued and the file in the other node is aligned to the
local changed one.

#### Get files sync state

Users can inspect the alignment of files through the offered swagger API

    curl -s  -X GET http://localhost:9876/api/1/state | jq .

#### Live split-brain

A split-brain can be the result of two different conditions:

1. entangled is started *after* a change in one or more of the monitored files or,
2. the nodes of the cluster cannot reach each other through the cluster network and a change in one or more of the monitored files occurs.

When in this situation, the tool delegates the resolution of the conflict to the operator of the cluster. 

If anyway there are further changes to the file and the two nodes are in reach, the change is propagated and the conflict is automatically resolved.

#### Force Node realignment

Users can force the alignment of one node to another throug the offered swagger API

    curl -X POST http://localhost:9876/api/1/force

This can be of help when we need to resolve a split-brain.