package main

import(
    // log package for error logging
    "log"
)

// Comment represents the structure of a comment
// with a Text field that can be populated from form data or JSON
type Comment struct {
    Text string `form:"text" json:"text" `
}

// main function initializes and starts the Fiber web server
func main(){
    // Create a new Fiber instance
    app:=fiber.New()
    
    // Create an API group with "/api/v1" prefix
    api:=app.Group("/api/v1")
    
    // Register POST endpoint for creating comments
    api.Post("/comments", createComment)
    
    // Start the server on port 3000
    app.Listen(":3000")
}

// createComment handles POST requests to create new comments
// Parameters:
//   - c *fiber.Ctx: Fiber context containing request/response data
// Returns:
//   - error: Any error that occurred during processing
func createComment(c *fiber.Ctx)error{
    // Initialize new Comment struct
    cmt:= new(Comment)
    
    // Parse request body into Comment struct
    if err:=c.BodyParser(cmt); err!=nil{
        // Log parsing error
        log.Println(err)
        // Return 400 Bad Request with error message
        c.Status(400).JSON(&fiber.Map{
            "success":false,
            "message": err,
        })
        return err
    }		
    
    // Convert comment to JSON bytes for queue processing
    cmtInBytes, err:=json.Marshal(cmt)
    // Push comment to Kafka queue
    PushCommentToQueue("comments", cmtInBytes)

    // Send success response with comment data
    err = c.JSON(&fiber.Map{
        "success":true,
        "message":"Comment pushed succcesfully",
        "comment":cmt,
    })
    if err != nil{
        // Handle JSON response error with 500 Internal Server Error
        c.Status(500).JSON(&fiber.Map{
            "success":false,
            "message": "Error creating product",
        })
        return err
    }
    return err
}