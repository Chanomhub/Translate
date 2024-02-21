﻿using System;
using System.IO;
using System.Text;
using XUnity.AutoTranslator.Plugin.ExtProtocol;

namespace ezTransXP.ExtProtocol
{
   class Program
   {
      static void Main( string[] args )
      {
         try
         {
            if( args.Length == 0 )
            {
               Console.WriteLine( "This program is automatically run during a game session if ezTrans is selected, and will be automatically stopped afterwards. This program cannot be run independently." );
               Console.WriteLine( "Press any key to exit." );
               Console.ReadKey();
               return;
            }
            var powerTranslatorPathPayload = args[ 0 ];
            var powerTranslatorPath = Encoding.UTF8.GetString( Convert.FromBase64String( powerTranslatorPathPayload ) );

            var datPath = Path.Combine( powerTranslatorPath, @"Dat" );  //initialize path

            var dllPath = Path.Combine( powerTranslatorPath, @"J2KEngineH.dll" );   //Ehnd engine path
            if( !File.Exists( dllPath ) )
            {
               dllPath = Path.Combine( powerTranslatorPath, @"J2KEngine.dll" );
            }


            using( var translator = new ezTransTranslationLibrary( dllPath ) )
            {
               if( !translator.InitEx( "CSUSER123455", datPath ) )  //initialize ezTrans
               {
                  throw new Exception( "J2K_InitializeEx Failed." );
               }
               using( var stdout = Console.OpenStandardOutput() )
               using( var writer = new StreamWriter( stdout ) )
               using( var stdin = Console.OpenStandardInput() )
               using( var reader = new StreamReader( stdin ) )
               {
                  writer.AutoFlush = true;

                  while( true )
                  {
                     var receivedPayload = reader.ReadLine();
                     if( string.IsNullOrEmpty( receivedPayload ) ) return;

                     var message = ExtProtocolConvert.Decode( receivedPayload ) as TranslationRequest;
                     if( message == null ) return;

                     var translatedTexts = new string[ message.UntranslatedTexts.Length ];
                     for( int i = 0; i < message.UntranslatedTexts.Length; i++ )
                     {
                        var untranslatedText = message.UntranslatedTexts[ i ];
                        var translatedText = translator.Translate( untranslatedText );
                        translatedTexts[ i ] = translatedText;
                     }

                     var response = new TranslationResponse
                     {
                        Id = message.Id,
                        TranslatedTexts = translatedTexts
                     };

                     var translatedPayload = ExtProtocolConvert.Encode( response );
                     writer.WriteLine( translatedPayload );
                  }
               }
            }
         }
         catch( Exception )
         {
            // "Graceful shutdown"
         }
      }
   }
}
