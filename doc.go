/*
Package prefixing is for resolving paths
*/
package prefixing

// http router tree impl seems to be the most performant
// but we will need to rebuild the tree each time for resolving routes
// while not ideal, this can be an RW lock easily to build/swap on.
