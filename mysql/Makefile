schema.sql :
	mysqldump -d -u $$USER -p $$DBNAME > schema.sql

clean :
	rm -f schema.sql

.PHONY : clean
