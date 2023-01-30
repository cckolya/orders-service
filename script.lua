
request = function()

   -- insert here key of values
   array = {"mnddphzlyr", "rghxrrvmel", "cwltgqxhnc", "yubppwsvnq", "dcoajmefku", "celkjecryh", "ercefeenfd", "yzgbogutws", "dcvwhkhade", "cqotcgstuo", "fvrpyyytjd", "gszcwyhscz", "yetaoxfmmm", "vhnuuskroj"}

   print(array[math.random(1,#array)])

   url_path = "/order/" .. array[math.random(1,#array)]
--    print(url_path)
   return wrk.format("GET", url_path)
end