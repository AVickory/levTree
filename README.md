LevTree

An eventually consistent tree database implemented on top of leveldb.

There's a lot of work left to go, I really don't take advantage of leveldb's
//sequential seek capabilities.  Plans are in the works.

Oh, and a proper README.  For now, if you'd like to check out the api you'll
have to use the comments or use godoc (Manually, I'll go get it hosted once
it's tested and usable).

tree.go

The Tree module sets up the data structures required to represent tree nodes in
memory.  The basic building block is the record which consists of a location
and some data (location basically being a key in the database).

Nodes contain a record with thier own location and the data that's meant
to be stored on that node.  They also contain a parent record which contains
thier parent's location and some data describing thier parent.  Additionally
they have a record for each of that node's children.

Within the provided API, note that a locateable can be either a record or a 
node.  This simplifies the api so that calling .UpdateNodeData on a node's
parent record is the same thing as calling it on the parent.

Within the database there are four kinds of nodes.  They are not different
types, again to simplify the API, but rather designate how thier children will
be bucketed relative to themselves.

Branch - A "normal" node.  Children of this kind of node will be in the same bucket
that this node is.  Most nodes should be Branches since other types will cause
the keys in the database to get progressively longer.  a single branch should
generally hold all of the data that you'll need on a given db access (within
reason; leveldb will run slowly if your entries get to long).

Tree - The root node of a single tree.  A tree is a node who's children will use the
tree's key as thier bucket.  They allow for grouping of related entries in the
db and allow for sequential access of all of thier descendants at once.  Trees
are generally meant to be the children of either other trees or forests.

Forest - A tree node that is attached to the root node of the db.  If you're
coming from a SQL background, then you might think about forests as your 
database's tables and a forest's child trees as sub tables.

Root - the bottom namespace of the database.  This node is just a place to hold
Meta data about your forests.  Since you can't directly name your forests, you
can instead put identifying data in the root's child metadata records.

dbFunnel.go

The DbFunnel module is meant to make writes take up less total time and to
ensure that consecutive updates of any document are 

One of the issues of goleveldb, is that if one transaction is open and you try
to open another then the you get an error instead of causing the thread to 
block.  To manage this problem I use a funnel.  This also allows for writes to
be periodically batch written to the db so that less total time is spent
writing and hence blocking reads.

One of the consequences of how this is implemented is that you should never
assume that an update that you just ran is actually available to you
through the provided read methods or is on the db.  All reads from the api go
to the database itself and bypass the funnel so that reads and writes don't
have to compete for access.  When an update is called for an node that is in
the funnel that update will be applied to that copy of the node in the funnel.

location.go

The Location module provides a bucketing system for namespacing keys.  id 
generation uses guuid V4, which is sufficient for my usecase, but other
options may be added later.