select p.iProdID, p.vName, pc.vName, pca.fRetailPrice 
	from product p 
        join prodcat pc on p.iPCatID = pc.iPCatID
        join product_color_assoc pca on p.iProdID = pca.iProdID
            where pca.fRetailPrice = 0.0;

select distinct pc.vName 
	from product p 
        join prodcat pc on p.iPCatID = pc.iPCatID
        join product_color_assoc pca on p.iProdID = pca.iProdID
            where 
                pca.fRetailPrice > 0.0 AND
                pca.iProdAttribID = 0;

select distinct p.vName, pc.vName, pca.fRetailPrice 
	from product p 
        join prodcat pc on p.iPCatID = pc.iPCatID
        join product_color_assoc pca on p.iProdID = pca.iProdID
            where 
                pca.iProdAttribID = 0;
