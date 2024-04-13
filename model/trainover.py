import os
import tensorflow as tf
from tensorflow.keras.preprocessing.image import ImageDataGenerator
from tensorflow.keras.models import load_model

# Load the pre-trained model
model = load_model('modelv4.h5')

# Set the directory where the new Pokemon images are stored
new_data_dir = 'PokemonData'

# Create the data generators for the new data
new_train_datagen = ImageDataGenerator(rescale=1./255)
new_train_generator = new_train_datagen.flow_from_directory(
    new_data_dir,
    target_size=(224, 224),
    batch_size=32,
    class_mode='categorical')

# Freeze the base model layers
for layer in model.layers[:-3]:
    layer.trainable = False

# Compile the model
optimizer = tf.keras.optimizers.Adam(learning_rate=0.001)
model.compile(optimizer=optimizer, loss='categorical_crossentropy', metrics=['accuracy'])

# Fine-tune the model
model.fit(new_train_generator,
          steps_per_epoch=None,
          epochs=5)

# Save the fine-tuned model
model.save('modelv41.h5')
