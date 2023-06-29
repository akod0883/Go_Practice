import csv
import random


def generate_sentence():
    length = random.randint(1, 20)  # Generate random sentence length between 1 and 20 words
    sentence = " ".join([random.choice(words) for _ in range(length)])
    return sentence

words = ["apple", "banana", "cherry", "date", "elderberry", "fig", "grape", "honeydew", "kiwi", "lemon"]

# Generate sentences
sentences = [generate_sentence() for _ in range(10)]

# Write sentences to CSV file
filename = 'sentences.csv'
with open(filename, 'w', newline='') as csvfile:
    writer = csv.writer(csvfile)
    writer.writerow(['Sentence'])

    for sentence in sentences:
        writer.writerow([sentence])

print(f"CSV file '{filename}' with one million sentences has been generated.")
