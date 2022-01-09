package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"context"
	"os"
	"github.com/joho/godotenv"
	"github.com/emvi/null"
	"golang.org/x/crypto/bcrypt"
)

// structs

// note: make sure the attributes in the struct are Capitalized
//		if not they won't be exported and cannot be accessed (kinda like private in Java),
//		as such, .BindJson will not be able to access attributes, 
//		causing an empty object ("{}") to be returned

type Credentials struct {
	Username string `json:username`
	Password string `json:password`
}

type User struct {
	Id int `json:"id"`
	Username string `json:"username"`
}

type Task struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Category_Id int `json:"category_id"`
	Category string `json:"category_name"`
	Deadline null.Time `json:"deadline"`
	Completed bool `json:"completed"`
	Created_at null.Time `json:"created_at"`
	Updated_at null.Time `json:"updated_at"`
}

type Category struct {
	Id int `json:"category_id"`
	Title string `json:"category_name"`
}

type CreateTaskParams struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Category_Id int `json:"category_id"`
	Deadline null.Time `json:"deadline"`
	User User `json:"user"`
}

type UpdateTaskParams struct {
	Id int `json:id`
	Title string `json:"title"`
	Description string `json:"description"`
	Category_Id int `json:"category_id"`
	Deadline null.Time `json:"deadline"`
}

type GetTaskByIdParams struct {
	Id int `json:"id"`
}

type CreateCategoryParams struct {
	Title string `json:"category_name"`
	User User `json:"user"` 
}

type QueryTasksParams struct {
	CategoryId int `json:"category_id"`
	CompletionStatus int `json:"completion_status"`
	User User `json:"user"` 
}

// CORS middleware
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

func main() {
	r := gin.Default()

	// allow CORS
	r.Use(CORSMiddleware());

	/* --------------------------------------------------------------- URL ENDPOINTS -------------- */
	// ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "Hello!")
	})

	// sign a new user up
	r.POST("/signup", func(c *gin.Context) {		
		var params Credentials;
		err := c.BindJSON(&params);
		assertJSONSuccess(c, err);

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), 8)

		details, err := signUp(params.Username, hashedPassword, c)

		if (err == nil) {
			c.JSON(200, details);
		}
	})

	// log in an existing user
	r.POST("/login", func(c *gin.Context) {
		var params Credentials;
		err := c.BindJSON(&params);
		assertJSONSuccess(c, err);		

		details, err := logIn(params.Username, params.Password, c)

		if (err == nil) {
			c.JSON(200, details);
		}
	})

	// get tasks by category id and completion status
	r.POST("/gettasks", func(c *gin.Context) {		
		var params QueryTasksParams
		err := c.BindJSON(&params);
		assertJSONSuccess(c, err)

		var taskList []Task = getTasks(params, c);

		c.JSON(200, taskList)
	})

	// get a specific task by id
	r.POST("/gettask", func(c *gin.Context) {		
		var params GetTaskByIdParams;
		err := c.BindJSON(&params);
		assertJSONSuccess(c, err);

		var t Task = getTask(params.Id, c);
		
		if (!c.IsAborted()) {
			c.JSON(200, t)
		}
	})

	// update a specific task by id
	r.POST("/updatetask", func(c *gin.Context) {		
		var params UpdateTaskParams
		err := c.BindJSON(&params);
		assertJSONSuccess(c, err);

		updateTask(params, c)

		if (!c.IsAborted()) {
			c.JSON(200, fmt.Sprintf("Successfully updated task with id: %v", params.Id))
		}
	})

	// mark a task as completed by id
	r.POST("/completetask", func(c *gin.Context) {		
		var params GetTaskByIdParams;
		err := c.BindJSON(&params)
		assertJSONSuccess(c, err);

		completeTask(params.Id, c);

		c.JSON(200, fmt.Sprintf("Successfully completed task with id: %v", params.Id))
	})

	// mark a task as incomplete by id
	r.POST("/incompletetask", func(c *gin.Context) {
		

		var params GetTaskByIdParams;
		err := c.BindJSON(&params)
		assertJSONSuccess(c, err);

		incompleteTask(params.Id, c);

		c.JSON(200, fmt.Sprintf("Successfully marked task as incomplete with id: %v", params.Id))
	})

	// deletes a task by id
	r.POST("/deletetask", func(c *gin.Context) {		
		var params GetTaskByIdParams;
		err := c.BindJSON(&params)
		assertJSONSuccess(c, err);

		deleteTask(params.Id, c);

		c.String(200, fmt.Sprintf("Successfully deleted task with id: %v", params.Id))
	})

	// adds a task
	r.POST("/addtask", func(c *gin.Context) {		
		var params CreateTaskParams
		err := c.BindJSON(&params)
		assertJSONSuccess(c, err);

		addTask(params, c)		
	})

	// create a category
	r.POST("/addcategory", func(c *gin.Context) {		
		var params CreateCategoryParams
		err := c.BindJSON(&params)
		assertJSONSuccess(c, err);

		addCategory(params, c)	
	})

	// gets a list of all categories
	r.POST("/allcategories", func(c *gin.Context) {		
		var params User;
		err := c.BindJSON(&params);
		assertJSONSuccess(c, err);

		var categoryList []Category = getAllCategories(params.Id, c);

		c.JSON(200, categoryList)
	})

	// start the server
	r.Run()
}

/* ----------------------------------------------------------------- DATABASE FUNCTIONS --------- */
/* Initialises and returns a connection to the database */
func connectDB(client *gin.Context) (c *pgx.Conn) {
	// load the .env file that contains postgresql connection details
	godotenv.Load(".env")

	// open a connection to the database
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	// check that the connection is successfully established
	assertDBSuccess(client, err);

	return conn;
}

/* Creates an account for a new user and returns the details of the user */
func signUp(username string, password []byte, client *gin.Context) (User, error) {
	c := connectDB(client)
	defer c.Close(context.Background())

	_, err := c.Exec(context.Background(), "INSERT INTO users (username, password) VALUES ($1, $2);", username, password);

	var user User

	// if a user with the same username already exists,
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "Unable to create a new user account: %v\n", err);

		// return HTTP Error 409: Conflict
		client.AbortWithStatusJSON(409, gin.H{"error": "Username already taken."});		

		return user, err;
	}

	err = c.QueryRow(context.Background(), "SELECT id, username FROM users WHERE username=$1", username).Scan(
		&user.Id,
		&user.Username,
	)
	assertDBOperationSuccess(client, err);

	return user, nil;
}

/* Logs in an existing-user */
func logIn(username string, password string, client *gin.Context) (User, error) {
	c := connectDB(client)
	defer c.Close(context.Background())

	var storedPassword string;
	var user User;

	err := c.QueryRow(context.Background(), "SELECT password FROM users WHERE username=$1;", username).Scan(&storedPassword);

	// check if the user exists
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "Unable to log user in: %v\n", err);

		// return HTTP Error 401: Unauthorised
		client.AbortWithStatusJSON(401, gin.H{"error": "Username doesn't exist."});		

		return user, err;
	}

	// check if the password is correct
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password));
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "Unable to log user in: %v\n", err);

		client.AbortWithStatusJSON(401, gin.H{"error": "Wrong password."});		

		return user, err;
	}

	// if credentials given are correct, return the user object
	err = c.QueryRow(context.Background(), "SELECT id, username FROM users WHERE username=$1", username).Scan(
		&user.Id,
		&user.Username,
	)
	assertDBOperationSuccess(client, err);

	return user, nil;
}

/* Returns an array of Tasks based on filtering criteria */
func getTasks(filterParams QueryTasksParams, client *gin.Context) ([]Task) {
	c := connectDB(client)
	defer c.Close(context.Background())

	var tasks pgx.Rows;
	var err error;

	if (filterParams.CategoryId == -1) {
		if (filterParams.CompletionStatus == 0) {
			/* get all tasks */
			tasks, err = c.Query(context.Background(), "SELECT * from public.get_all_tasks($1) ORDER BY id;", filterParams.User.Id);
		} else if (filterParams.CompletionStatus == 1) {
			/* get all completed tasks */
			tasks, err = c.Query(context.Background(), "SELECT * from public.get_completed_tasks($1) ORDER BY id;", filterParams.User.Id);
		} else {
			/* get all incomplete tasks */
			tasks, err = c.Query(context.Background(), "SELECT * from public.get_incomplete_tasks($1) ORDER BY id;", filterParams.User.Id);
		}
	} else {
		if (filterParams.CompletionStatus == 0) {
			/* get all tasks tagged with the category */
			tasks, err = c.Query(context.Background(), "SELECT * from public.get_tasks_in_category($1) ORDER BY id;", filterParams.CategoryId);
		} else if (filterParams.CompletionStatus == 1) {
			/* get all completed tasks tagged with the category */
			tasks, err = c.Query(context.Background(), "SELECT * from public.get_tasks_in_category($1) WHERE completed='t' ORDER BY id;", filterParams.CategoryId);
		} else {
			/* get all incomplete tasks tagged with the category */
			tasks, err = c.Query(context.Background(), "SELECT * from public.get_tasks_in_category($1) WHERE completed='f' ORDER BY id;", filterParams.CategoryId);
		}
	}

	var taskSlice []Task
	for tasks.Next() {
		var t Task
		err = tasks.Scan(
			&t.Id, 
			&t.Title,
			&t.Description,
			&t.Category_Id,
			&t.Category,
			&t.Deadline,
			&t.Completed,
			&t.Created_at,
			&t.Updated_at,	
		)
		assertDBOperationSuccess(client, err);
		taskSlice = append(taskSlice, t)
	}

	return taskSlice;
}

/* Return a Task by its id */
func getTask(task_id int, client *gin.Context) (Task) {
	c := connectDB(client)
	defer c.Close(context.Background())

	var t Task

	err := c.QueryRow(context.Background(), "SELECT * FROM public.get_task_by_id($1);", task_id).Scan(
		&t.Id, 
		&t.Title,
		&t.Description,
		&t.Category_Id,
		&t.Category,
		&t.Deadline,
		&t.Completed,
		&t.Created_at,
		&t.Updated_at,		
	)
	assertDBOperationSuccess(client, err);

	return t;
}

/* Update a Task by its id */
func updateTask(t UpdateTaskParams, client *gin.Context) {
	c := connectDB(client)
	defer c.Close(context.Background())

	_, err := c.Exec(context.Background(), "UPDATE tasks SET category_id=$1, title=$2, description=$3, deadline=$4 WHERE id=$5;", t.Category_Id, t.Title, t.Description, t.Deadline, t.Id)
	assertDBOperationSuccess(client, err);
}

/* Mark a Task as completed by its id */
func completeTask(id int, client *gin.Context) {
	c := connectDB(client)
	defer c.Close(context.Background())

	_, err := c.Exec(context.Background(), "UPDATE tasks SET completed='t' WHERE id=$1;", id);
	assertDBOperationSuccess(client, err);
}

/* Mark a previously completed task as incomplete by its id */
func incompleteTask(id int, client *gin.Context) {
	c := connectDB(client)
	defer c.Close(context.Background())

	_, err := c.Exec(context.Background(), "UPDATE tasks SET completed='f' WHERE id=$1;", id);
	assertDBOperationSuccess(client, err);
}

/* Deletes a Task in the database with the corresponding id */
func deleteTask(id int, client *gin.Context) {
	c := connectDB(client)
	defer c.Close(context.Background())

	// use Exec to execute a query that does not return a result set
	commandTag, err := c.Exec(context.Background(), "DELETE FROM tasks where id=$1;", id)
	assertDBOperationSuccess(client, err);
	if commandTag.RowsAffected() != 1 {
		fmt.Fprintf(os.Stderr, "No row found to delete\n")
		client.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return;
	}
}

/* Adds a Task to the database */
func addTask(params CreateTaskParams, client *gin.Context) {
	c := connectDB(client)
	defer c.Close(context.Background())

	commandTag, err := c.Exec(context.Background(), "INSERT INTO tasks (category_id, title, description, deadline, user_id, completed) VALUES ($1, $2, $3, $4, $5, FALSE);", params.Category_Id, params.Title, params.Description, params.Deadline, params.User.Id)
	assertDBOperationSuccess(client, err);
	if commandTag.RowsAffected() != 1 {
		fmt.Fprintf(os.Stderr, "Task not added to db\n")
		client.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return;
	}
}

/* Adds a Category to the database */
func addCategory(params CreateCategoryParams, client *gin.Context) {
	c := connectDB(client)
	defer c.Close(context.Background())

	commandTag, err := c.Exec(context.Background(), "INSERT INTO categories (user_id, title) VALUES ($1, $2);", params.User.Id, params.Title)
	assertDBOperationSuccess(client, err);
	if commandTag.RowsAffected() != 1 {
		fmt.Fprintf(os.Stderr, "Category not added to db\n")
		client.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return;
	}
}

/* Returns a list of categories with their associated primary-keys */
func getAllCategories(user_id int, client *gin.Context) ([]Category) {
	c := connectDB(client)
	defer c.Close(context.Background())

	categories, err := c.Query(context.Background(), "SELECT id, title from categories WHERE user_id=$1 ORDER BY id;", user_id)
	assertDBOperationSuccess(client, err);
	defer categories.Close();

	var categorySlice []Category
	for categories.Next() {
		var cat Category
		err = categories.Scan(
			&cat.Id,
			&cat.Title,
		)
		assertDBOperationSuccess(client, err);
		categorySlice = append(categorySlice, cat)
	}

	return categorySlice;
}

/* ------------------------------------------------------------ HELPER FUNCTIONS --------------------- */
// checks if there is an error connecting to the database,
//		if so, returns an error message to the client and cancels the context of the caller
func assertDBSuccess(client *gin.Context, e error) {
	if (e != nil) {
		// print error message on server side so that its visible in the server logs
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", e);

		// return http code of 500 to the client, which stands for "Internal Server Error"
		client.AbortWithStatusJSON(500, gin.H{"error": e.Error()});

		// halts execution of remaining functions to not do unnecessary work
		
	}
}

// checks if there is an error performing the specified request on the database,
//		if so, returns an error message to the client and cancels the context of the caller
func assertDBOperationSuccess(client *gin.Context, e error) {
	if (e != nil) {
		// print error message on server side so that its visible in the server logs
		fmt.Fprintf(os.Stderr, "Unable to perform the requested action: %v\n", e);

		// return http code of 500 to the client, which stands for "Internal Server Error"
		//		and halts execution of remaining functions to not do unnecessary work
		client.AbortWithStatusJSON(500, gin.H{"error": e.Error()});		
	}
}

// checks if there is an error connecting to the parsing JSON body,
//		if so, returns an error message to the client and stops execution of any remaining function-calls
func assertJSONSuccess(client *gin.Context, e error) {
	if (e != nil) {
		// print error message on server side so that its visible in the server logs
		fmt.Fprintf(os.Stderr, "Unable to parse JSON body: %v\n", e);

		// return http code of 406 to the client, which stands for "Not Acceptable" and
		// 		halts execution of remaining functions to not do unnecessary work		
		client.AbortWithStatusJSON(406, gin.H{"error": e.Error()});		
	}
}