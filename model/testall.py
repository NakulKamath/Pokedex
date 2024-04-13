import os
import tensorflow as tf
from tensorflow.keras.preprocessing.image import ImageDataGenerator

# Load the trained model
model = tf.keras.models.load_model('modelv4.h5')

# Define the directory for the validation dataset
validation_dir = 'dataset'

# Create a data generator for validation data
validation_datagen = ImageDataGenerator(rescale=1./255)

validation_generator = validation_datagen.flow_from_directory(
    validation_dir,
    target_size=(224, 224),
    batch_size=32,
    class_mode='categorical',
    shuffle=False)

# Evaluate the model on the validation dataset
loss, accuracy = model.evaluate(validation_generator, verbose=1)

print(f'Validation accuracy: {accuracy * 100:.2f}%')
