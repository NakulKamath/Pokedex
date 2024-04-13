import os
import pickle
from tensorflow.keras.preprocessing.image import load_img, img_to_array
import numpy as np

# Load the pickled model
with open('model.pkl', 'rb') as f:
    model = pickle.load(f)

# Define the target image size
target_size = (224, 224)

# Load and preprocess the input image
input_image = load_img('test10.png', target_size=target_size)
input_array = img_to_array(input_image) / 255.0
input_array = np.expand_dims(input_array, axis=0)

# Make the prediction
predictions = model.predict(input_array)

# Get the list of Pokemon names
with open('pokemon_names.pkl', 'rb') as f:
    pokemon_names = pickle.load(f)

# Get the top 3 predictions
top_3_indices = np.argsort(predictions[0])[-3:][::-1]
top_3_probabilities = [predictions[0][i] for i in top_3_indices]
top_3_pokemon = [pokemon_names[i] for i in top_3_indices]

# Print the results
print("Top 3 Predicted Pokemon:")
for i in range(3):
    print(f"{i+1}. Pokemon: {top_3_pokemon[i]}, Probability: {top_3_probabilities[i]:.2f}")