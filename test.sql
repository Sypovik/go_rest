SELECT n1.id, n1.title, n2.count
FROM note AS n1 
LEFT join (
    SELECT  title, content, count(title)
    FROM note GROUP BY title, content
    ) AS n2 
    ON n1.title=n2.title;



SELECT  title, content, count(title)
    FROM note GROUP BY title, content

SELECT id, title 
FROM note