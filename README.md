# GoSentry

Go service that monitors security logs for anomalies and sends notifications




# Configuration

Configurations are found in the c.yml file.

Option            | Explanation                                        | type 
----------------- | -------------------------------------------------- | ---------
flaggedPhrases    | List of phrases to look for in the audit files.    | []string
yearsBack         | Years back to look in the audit files.             | int      
monthsBack        | Months back to look in the audit files.            | int      
daysBack          | Days back to look in the audit files.              | int      
filesToWatch      | List of files to watch. (Full Paths).              | []string 
scanEvery         | How long between scans (Example 5h or 10d or 11M). | string   
outputDir         | What directory to output the flagged lines.        | string   
timeRegexPatterns | Regular Expressions used to parse the audit files. | []string 