package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const SIZE = 5

type Node struct {
	Val   string
	Left  *Node
	Right *Node
}

type Queue struct {
	Head   *Node
	Tail   *Node
	Length int
}

type Cache struct {
	Queue Queue
	Hash  Hash
}

type Hash map[string]*Node

func NewCache() Cache {
	return Cache{Queue: NewQueue(), Hash: Hash{}}
}

func NewQueue() Queue {
	head := &Node{}
	tail := &Node{}

	head.Right = tail
	tail.Left = head

	return Queue{Head: head, Tail: tail}
}

func (c *Cache) Check(str string) {
	node := &Node{}

	if val, ok := c.Hash[str]; ok {
		node = c.Remove(val)
	} else {
		node = &Node{Val: str}
	}
	c.Add(node)
	c.Hash[str] = node
}

func (c *Cache) Remove(n *Node) *Node {
	fmt.Printf("remove: %s\n", n.Val)
	left := n.Left
	right := n.Right

	left.Right = right
	right.Left = left
	c.Queue.Length -= 1
	delete(c.Hash, n.Val)
	return n
}

func (c *Cache) Add(n *Node) {
	fmt.Printf("add: %s\n", n.Val)
	tmp := c.Queue.Head.Right

	c.Queue.Head.Right = n
	n.Left = c.Queue.Head
	n.Right = tmp
	tmp.Left = n

	c.Queue.Length++
	if c.Queue.Length > SIZE {
		c.Remove(c.Queue.Tail.Left)
	}
}

func (c *Cache) Display() {
	c.Queue.Display()
}

func (q *Queue) Display() {
	node := q.Head.Right
	fmt.Printf("%d - [", q.Length)
	for i := 0; i < q.Length; i++ {
		fmt.Printf("{%s}", node.Val)
		if i < q.Length-1 {
			fmt.Printf("<-->")
		}
		node = node.Right
	}
	fmt.Println("]")
}

func (c *Cache) ToList() []string {
	var list []string
	node := c.Queue.Head.Right
	for i := 0; i < c.Queue.Length; i++ {
		list = append(list, node.Val)
		node = node.Right
	}
	return list
}

func main() {
	cache := NewCache()
	for _, word := range []string{"parrot", "avocado", "dragonfruit", "tree", "potato", "tomato", "tree", "dog"} {
		cache.Check(word)
	}

	router := gin.Default()
	router.Static("/static", "./static")
	router.LoadHTMLFiles("index.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/cache", func(c *gin.Context) {
		c.JSON(http.StatusOK, cache.ToList())
	})

	router.POST("/cache", func(c *gin.Context) {
		var req struct {
			Value string `json:"value"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		cache.Check(req.Value)
		c.JSON(http.StatusOK, cache.ToList())
	})

	router.DELETE("/cache", func(c *gin.Context) {
		var req struct {
			Value string `json:"value"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if node, ok := cache.Hash[req.Value]; ok {
			cache.Remove(node)
		}
		c.JSON(http.StatusOK, cache.ToList())
	})

	fmt.Println("Server is running on http://localhost:8080")
	router.Run(":8080")
}
