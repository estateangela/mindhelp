class Resource {
  final String id;
  final String name;
  final String type;
  final String address;
  final String phone;
  final String website;
  final String description;
  final double latitude;
  final double longitude;

  Resource({
    required this.id,
    required this.name,
    required this.type,
    required this.address,
    required this.phone,
    required this.website,
    required this.description,
    required this.latitude,
    required this.longitude,
  });

  factory Resource.fromJson(Map<String, dynamic> json) {
    return Resource(
      id: json['id'] as String,
      name: json['name'] as String,
      type: json['type'] as String,
      address: json['address'] as String,
      phone: json['phone'] as String,
      website: json['website'] as String,
      description: json['description'] as String,
      latitude: json['latitude'] as double,
      longitude: json['longitude'] as double,
    );
  }
}
