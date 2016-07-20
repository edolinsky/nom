# nom
## A collection of food-related projects

### mineChefDb.go
Accesses chefdb.com's restaurant pages, extracts information, and saves that information in a MySQL database.  
_**Please do not use.**  
This was created as an exercise in web scraping, and I have since deleted all extracted data. Check out [the source](http://chefdb.com/pl) if you would like to access information from [the Chef and Restaurant Database](http://chefdb.com)._  
All due credit to Brian and George for creating and maintaining the Chef and Restaurant Database.

#### Requires:
1. MySQL database with create table syntax as follows:  
    ```
    	CREATE TABLE `restaurants` (
    	  `id` int(11) NOT NULL AUTO_INCREMENT,
    	  `name` varchar(45) DEFAULT NULL,
    	  `rating` double DEFAULT NULL,
    	  `address` varchar(45) DEFAULT NULL,
    	  `city` varchar(45) DEFAULT NULL,
    	  `region` varchar(45) DEFAULT NULL,
    	  `url` varchar(45) DEFAULT NULL,
    	  `phone` varchar(45) DEFAULT NULL,
    	  PRIMARY KEY (`id`),
    	  UNIQUE KEY `id_UNIQUE` (`id`)
    	) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;
    ```
2. `github.com/go-sql-driver/mysql` package
