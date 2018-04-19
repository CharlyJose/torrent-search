
require 'rubygems'
require 'net/http'
require 'json'
require 'rainbow'
require 'terminal-table'	


tab = '     '
miniTab = '  '

puts "1. Movie\n2. TV Series"
user_choice = gets.chomp


case user_choice	
when "1"
    puts "Movie you want to search for?"
    user_movie_name = gets.chomp
    movie_name = user_movie_name.downcase.tr(" ", "-")
    url = "https://yts.am/api/v2/list_movies.json?query_term="+ movie_name + '"'
    uri = URI(url)
    response = Net::HTTP.get(uri)
    parsed = JSON.parse(response)
    query_response_status = parsed["status"]

    # If Search  successful and return status
    if query_response_status=="ok"  

        puts
        puts
        
        movie_cnt = 0
        movie_cnt = parsed["data"]["movie_count"]
        puts "Number of Movies Available : #{movie_cnt}" 
        puts
        
        movies_array = parsed["data"]["movies"]
            
        if  movie_cnt>0
            count=0
            parsed["data"]["movies"].each do |movies|  
                title=movies["title_long"]
                rating=movies["rating"]
                count+=1
                puts
                puts Rainbow("#{miniTab}##{count} - #{title} ").bg(:white).black + Rainbow(" Rating : #{rating} ").bg(:blue).white.bright
                puts
            
                movies["torrents"].each do |torrents|
                    peers=torrents["peers"].to_i
                    seeds=torrents["seeds"].to_i
                    
                    puts "Quality: "+ torrents["quality"]+" Size: "+ torrents["size"]  + Rainbow(" Seeders: #{seeds}").green  + Rainbow(" Peers: #{peers} ").red + Rainbow("  Torrent: #{torrents["url"]}").yellow
                end
        end
        
        elsif movie_cnt==0
            puts "Sorry !! No Movies found with  '#{user_movie_name}' as name.. "
        end

        puts
        puts

    elsif query_response_status=="error"
        puts "Error!! Your Search query could not be processed"
    else 
        puts "Unknown Error. Please report this"
end

when "2"
    puts "Series you want to search for?"
    user_series_name = gets.chomp
    series_name = user_series_name.downcase.tr(" ", "-")
    API_KEY = "YOU KEY HERE"
    url = "http://www.omdbapi.com/?apikey=" + API_KEY + "&t=" + series_name
    uri = URI(url)
    response = Net::HTTP.get(uri)
    parsed = JSON.parse(response)
    query_response_status = parsed["Response"]
    rating = parsed["imdbRating"]
    total_seasons = parsed["totalSeasons"]
    imdb_id = parsed["imdbID"]
    
    puts
    puts Rainbow(" Imdb ID: #{imdb_id}#{tab}").bg(:cyan).white.bright + Rainbow(" Rating : #{rating}#{tab}").bg(:blue).white.bright + Rainbow(" Total Seasons : #{total_seasons}#{tab}").bg(:yellow).white.bright
    

    #If Search  successful and return status
    if query_response_status=="True"  

        puts

        #Trimming first two charaters. EZTV API demands imdb id without first two charaters
        imdb_id[0] = ''
        imdb_id[0] = ''

        url = "https://eztv.ag/api/get-torrents?&imdb_id=" + imdb_id
        uri = URI(url)
        response2 = Net::HTTP.get(uri)
        parsed2 = JSON.parse(response2)
        
        puts

        torrent_cnt = 0
        torrent_cnt = parsed2['torrents'].length
        puts Rainbow("#{miniTab}Showing #{torrent_cnt} torrents#{tab}").bg(:white).red.bright
        
        puts

        if  torrent_cnt>0
            count = 0
            parsed2["torrents"].each do |series|
                count+=1
                title = series["title"]
                filename = series["filename"]
                video_resolution = 'Not Provided'
                video_codec = ''
                video_file_format = 'Not provided'

                # Video resolution
                if( title =~ /720p(.*)/ )
                    video_resolution = '720p'
                end
                if( title =~ /1080p(.*)/ )
                    video_resolution = '1080p'
                end

                # Video codec
                if( title =~ /WEBRip(.*)/ )
                    video_codec = video_codec + ' WEBRip'
                end
                if( title =~ /HDTV(.*)/ )
                    video_codec = video_codec + ' HDTV'
                end
                if( title =~ /Xvid(.*)/ )
                    video_codec = video_codec + ' Xvid'
                end
                if( title =~ /x264(.*)/ )
                    video_codec = video_codec + ' x264'
                end
                if( title =~ /Dvix(.*)/ )
                    video_codec = video_codec + ' Dvix'
                end
                if( title =~ /CamRip(.*)/ )
                    video_codec = video_codec + ' CamRip'
                end
                if( title =~ /DVDRip(.*)/ )
                    video_codec = video_codec + ' DVDRip'
                end
                if( title =~ /HDRip(.*)/ )
                    video_codec = video_codec + ' HDRip'
                end
                if( title =~ /WEBDL(.*)/ )
                    video_codec = video_codec + ' WEBDL'
                end
                if( title =~ /BRRip(.*)/ )
                    video_codec = video_codec + ' BRRip'
                end

                # Video file format
                if( filename =~ /webm(.*)/ )
                    video_file_format = 'webm'
                end
                if( filename =~ /mkv(.*)/ )
                    video_file_format = 'mkv'
                end
                if( filename =~ /avi(.*)/ )
                    video_file_format = 'avi'
                end
                if( filename =~ /m4v(.*)/ )
                    video_file_format = 'm4v'
                end
                if( filename =~ /3gp(.*)/ )
                    video_file_format = '3gp'
                end
                if( filename =~ /mp4(.*)/ )
                    video_file_format = 'mp4'
                end

                puts
                puts

                puts miniTab + Rainbow(" ##{count} - #{title} ").bg(:white).black
                puts
                puts miniTab + Rainbow('Seeders: ' + series['seeds'].to_s ).bg(:green).white.bright  + tab + Rainbow('Peers: ' + series['peers'].to_s ).bg(:red).white.bright + tab +Rainbow('Size: ' + ( series['size_bytes'].to_f/1073741824 ).round(4).to_s + ' GB' ).bg(:blue).white.bright
                puts miniTab + Rainbow('Resolution: '+ video_resolution).green.bright + tab + Rainbow('Codec:' + video_codec).red.bright + tab + Rainbow('Format: ' + video_file_format).blue.bright
                puts miniTab + Rainbow('Episode URL: ').white.bright + series['episode_url']
                puts miniTab + Rainbow('Torrent URL: ').white.bright + series['torrent_url']
                puts miniTab + Rainbow('Magnet URL: ').white.bright + Rainbow(series['magnet_url']).yellow
                puts miniTab + Rainbow('Torrent Hash: ').white.bright + series['hash']
                puts miniTab + Rainbow('Release date: ').white.bright + Time.at( series['date_released_unix'] ).to_s
            end

        elsif torrent_cnt==0
            puts "Sorry !! No Movies found with  '#{user_series_name}' as name.. "
        end

        puts
        puts


    elsif query_response_status=="False"
        query_error_report = parsed["Error"]

        puts "Error!! Your Search query could not be processed"
        puts "Error: " + Rainbow(query_error_report).red
    else 
        puts "Unknown Error. Please report this"
            
    end
    
end