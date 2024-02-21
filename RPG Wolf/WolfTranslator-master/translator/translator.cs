using GoogleTranslateFreeApi;
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Net;
using System.Text;
using System.Text.RegularExpressions;
using System.Threading;
using System.Threading.Tasks;
using static System.Net.Mime.MediaTypeNames;

namespace translator
{
    class translator
    {
        public static HttpListener listener;
        public static string url = "http://localhost:5731/";
        public static GoogleTranslator _translator = new GoogleTranslator();
        public static Dictionary<string, string> lookup = new Dictionary<string, string>();

        public static async Task HandleIncomingConnections()
        {
            bool runServer = true;

            while (runServer)
            {
                HttpListenerContext ctx = await listener.GetContextAsync();

                HttpListenerRequest req = ctx.Request;
                HttpListenerResponse resp = ctx.Response;
                Console.WriteLine(req.Url.ToString());

                Console.WriteLine(req.Url.ToString());
                Console.WriteLine(req.HttpMethod);
                Console.WriteLine(req.UserHostName);
                Console.WriteLine(req.UserAgent);
                Console.WriteLine();


                var body = req.InputStream;
                System.IO.StreamReader reader = new System.IO.StreamReader(body, Encoding.UTF8);

                var input = reader.ReadToEnd();
                Console.WriteLine(input);
                string output;
                if (new Regex(@"^[ -~\n\r]*$").IsMatch(input) || input.Trim() == "")
                {
                    output = input;
                } else
                {
                    output = await Translate(input);
                }
                byte[] data = Encoding.GetEncoding("ks_c_5601-1987").GetBytes(output);
                Console.WriteLine(output);
                resp.ContentEncoding = Encoding.GetEncoding("ks_c_5601-1987");
                resp.StatusCode = 200;
                await resp.OutputStream.WriteAsync(data, 0, data.Length);
                resp.Close();
                
            }
        }

        private static async Task<string> Translate(string buf)
        {
            int retry = 10;
            var buf_escaped = Regex.Escape(buf);
            if (lookup.ContainsKey(buf_escaped))
            {
                Console.WriteLine("Cache hit! " + buf_escaped);
                return Regex.Unescape(lookup[buf_escaped]);
            }
            while (retry != 0)
            {
                try
                {
                    Language from = Language.Auto;
                    Language to = GoogleTranslator.GetLanguageByName("Korean");

                    TranslationResult result = await _translator.TranslateLiteAsync(buf, from, to);
                    string resultMerged = result.MergedTranslation;
                    var translated_escpaed = Regex.Escape(resultMerged);
                    lookup[buf_escaped] = translated_escpaed;
                    Save();
                    return resultMerged;
                }
                catch (Exception error)
                {
                    _translator = new GoogleTranslator();
                    Console.WriteLine(error);
                    retry--;
                }
            }
            return "ERROR";
        }
        
        static void Save()
        {
            using (StreamWriter file = new StreamWriter("translation.txt"))
                foreach (var entry in lookup)
                    file.WriteLine("{0}\n{1}", entry.Key, entry.Value);
        }

        static void Load()
        {
            using (StreamReader file = new StreamReader("translation.txt"))
            {
                string line;
                while ((line = file.ReadLine()) != null)
                {
                    var key = line;
                    var value = file.ReadLine();
                    lookup[key] = value;
                }
            }
        }

        static int Main(string[] args)
        {
            if (File.Exists("translation.txt"))
            {
                Load();
            }
            listener = new HttpListener();
            listener.Prefixes.Add(url);
            listener.Start();
            Console.WriteLine("Listening for connections on {0}", url);
            HandleIncomingConnections().GetAwaiter().GetResult();
            listener.Close();
            return 0;
        }
    }
}
