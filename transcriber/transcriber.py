from vosk import Model, KaldiRecognizer
import sys, os
import subprocess
import json
import wave
import re

model = Model('vosk-model-small-ru-0.22')
inputfile = sys.argv[1]
wavfile = inputfile + '.wav'

print(inputfile)
print(os.getcwd())
subprocess.call(f"ffmpeg -y -i {inputfile} -ar 48000 -ac 1 -f wav {wavfile}")

wf = wave.open(wavfile, "rb")
rcgn_fr = wf.getframerate() * wf.getnchannels()
rec = KaldiRecognizer(model, rcgn_fr)
result = ''
last_n = False

read_block_size = wf.getnframes()

while True: 
    data = wf.readframes(read_block_size)
    if len(data) == 0:
        break

    if rec.AcceptWaveform(data):
        res = json.loads(rec.Result())
        
        if res['text'] != '':
            result += f" {res['text']}"
            if read_block_size < 200000:
                print(res['text'] + " \n")
            
            last_n = False
        elif not last_n:
            result += '\n'
            last_n = True

res = json.loads(rec.FinalResult())
result += f" {res['text']}"

output_file = open("output.txt", 'w')
print('\n'.join(line.strip() for line in re.findall(r'.{1,150}(?:\s+|$)', result)), file=output_file)
output_file.close()
