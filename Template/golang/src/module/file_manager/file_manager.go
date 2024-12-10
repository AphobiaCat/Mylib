package file_manager

import(
	"mylib/src/public"
	"os"
)

func File_Read(file_path string) string{
	data, err := os.ReadFile(file_path)
	if err != nil {
		public.DBG_ERR("Error reading file:", err)
		return ""
	}

	return string(data)
}

func File_Append(write_data string, file_path string, batchSize int){
	// create file if no exist
	file, err := os.OpenFile(file_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		public.DBG_ERR("Error opening file:", err)
		return
	}
	defer file.Close()

	for i := 0; i < len(write_data); i += batchSize {
		end := i + batchSize
		if end > len(write_data) {
			end = len(write_data)
		}
		batch := write_data[i:end]
		_, err := file.Write([]byte(batch))
		if err != nil {
			public.DBG_ERR("Error write file:", err)
			return
		}
	}
	return
}


