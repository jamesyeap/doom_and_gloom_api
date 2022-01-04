-- get all tasks
CREATE OR REPLACE FUNCTION public.get_all_tasks(Specified_User_Id INT)
	RETURNS TABLE 
		(
			id INT,
			title VARCHAR(255),
			description TEXT,
			category_id INT,
			category TEXT,
			deadline TIMESTAMP,
			completed BOOLEAN,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)
	language plpgsql
AS
$$
BEGIN
	RETURN QUERY
		SELECT 
			tasks.id,
			tasks.title,
			tasks.description,
			categories.id,
			categories.title,
			tasks.deadline,
			tasks.completed,
			tasks.created_at,
			tasks.updated_at
		FROM
			public.tasks
				INNER JOIN public.categories ON public.tasks.category_id=public.categories.id
		WHERE
			tasks.user_id=Specified_User_Id;
END
$$;

-- get task by id
CREATE OR REPLACE FUNCTION public.get_task_by_id(Specified_User_Id INT, Specified_Task_Id INT)
	RETURNS TABLE 
		(
			id INT,
			title VARCHAR(255),
			description TEXT,
			category_id INT,
			category TEXT,
			deadline TIMESTAMP,
			completed BOOLEAN,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)
	language plpgsql
AS
$$
BEGIN
	RETURN QUERY
		SELECT 
			tasks.id,
			tasks.title,
			tasks.description,
			categories.id,
			categories.title,
			tasks.deadline,
			tasks.completed,
			tasks.created_at,
			tasks.updated_at
		FROM
			public.tasks
				INNER JOIN public.categories ON public.tasks.category_id=public.categories.id
		WHERE
			tasks.user_id=Specified_User_Id AND
			tasks.id=Specified_Task_Id;
END
$$;

-- get all completed tasks
CREATE OR REPLACE FUNCTION public.get_completed_tasks(Specified_User_Id INT)
	RETURNS TABLE
		(
			id INT,
			title VARCHAR(255),
			description TEXT,
			category_id INT,
			category TEXT,
			deadline TIMESTAMP,
			completed BOOLEAN,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)
	language plpgsql
AS
$$
BEGIN
	RETURN QUERY
		SELECT 
			tasks.id,
			tasks.title,
			tasks.description,
			categories.id,
			categories.title,
			tasks.deadline,
			tasks.completed,
			tasks.created_at,
			tasks.updated_at
		FROM
			public.tasks
				INNER JOIN public.categories ON public.tasks.category_id=public.categories.id

		WHERE
			tasks.user_id = Specified_User_Id AND
			tasks.completed = 't';
END
$$;

-- get all outstanding tasks
CREATE OR REPLACE FUNCTION public.get_incomplete_tasks(Specified_User_Id INT)
	RETURNS TABLE
		(
			id INT,
			title VARCHAR(255),
			description TEXT,
			category_id INT,
			category TEXT,
			deadline TIMESTAMP,
			completed BOOLEAN,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)
	language plpgsql
AS
$$
BEGIN
	RETURN QUERY
		SELECT 
			tasks.id,
			tasks.title,
			tasks.description,
			categories.id,
			categories.title,
			tasks.deadline,
			tasks.completed,
			tasks.created_at,
			tasks.updated_at
		FROM
			public.tasks
				INNER JOIN public.categories ON public.tasks.category_id=public.categories.id

		WHERE
			tasks.user_id = Specified_User_Id AND
			tasks.completed = 'f';
END
$$;

-- get tasks by category id
CREATE OR REPLACE FUNCTION public.get_tasks_in_category(Specified_User_Id INT, Specified_Category_Id INT)
	RETURNS TABLE
		(
			id INT,
			title VARCHAR(255),
			description TEXT,
			category_id INT,
			category TEXT,
			deadline TIMESTAMP,
			completed BOOLEAN,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)
	language plpgsql
AS
$$
BEGIN
	RETURN QUERY
		SELECT 
			tasks.id,
			tasks.title,
			tasks.description,
			categories.id,
			categories.title,
			tasks.deadline,
			tasks.completed,
			tasks.created_at,
			tasks.updated_at
		FROM
			public.tasks
				INNER JOIN public.categories ON public.tasks.category_id=public.categories.id

		WHERE
			tasks.user_id = Specified_User_Id AND
			categories.id = Specified_Category_Id;
END
$$;
